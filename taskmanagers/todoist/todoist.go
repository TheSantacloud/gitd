package todoist

import (
	"encoding/json"
	"fmt"
	"log"
	"mgtd/adapters"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TodoistAdapter struct {
	endpointURL string
	httpClient  *http.Client
	authToken   string
	settings    adapters.Settings
}

func (t *TodoistAdapter) Initialize(settings adapters.Settings) error {
	t.endpointURL = "https://api.todoist.com/sync/v9/sync"
	t.httpClient = &http.Client{
		Timeout: 15 * time.Second, // Todoist default timeout
	}
	t.settings = settings

	authToken, err := GenerateAccessToken(settings.Todoist)
	if err != nil {
		return err
	}
	t.authToken = authToken

	return nil
}

func (t *TodoistAdapter) FetchTasks() ([]adapters.Task, error) {
	data := url.Values{}
	data.Set("sync_token", "*")
	data.Set("resource_types", "[\"items\",\"projects\",\"notes\",\"labels\",\"sections\"]")
	req, err := http.NewRequest(http.MethodPost, t.endpointURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+t.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := t.httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result TodoistSyncResponse
	decoder := json.NewDecoder(res.Body)

	if err := decoder.Decode(&result); err != nil {
		fmt.Println("error decoding response", err)
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("Status code error: %d for url: %s", res.StatusCode, req.URL)
	}

	tasks := result.ToTasks()
	return tasks, nil
}

func (t *TodoistAdapter) UpdateTasks(actions *[]adapters.TaskAction) error {
	syncResponse := []SyncResponseItem{}

	prepareCompletedSync(actions, &syncResponse)
	prepareDeletedSync(actions, &syncResponse)
	prepareDeferredSync(actions, &syncResponse)
	prepareRevalidateSync(actions, &syncResponse)

	jsonBytes, err := json.MarshalIndent(syncResponse, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}

	fmt.Println(string(jsonBytes))
	return nil
}

func prepareCompletedSync(actions *[]adapters.TaskAction, syncResponse *[]SyncResponseItem) {
	tasksToComplete := getAllTasksWithAction(actions, adapters.ActionComplete)
	if len(*tasksToComplete) == 0 {
		return
	}
	idsToComplete := extractIds(tasksToComplete)
	*syncResponse = append(*syncResponse, SyncResponseItem{
		Type: "item_complete",
		Uuid: uuid.New().String(),
		Args: &SyncResponseArgs{
			Ids: idsToComplete,
		},
	})
}

func prepareDeletedSync(actions *[]adapters.TaskAction, syncResponse *[]SyncResponseItem) {
	tasksToDelete := getAllTasksWithAction(actions, adapters.ActionDelete)
	if len(*tasksToDelete) == 0 {
		return
	}
	idsToDelete := extractIds(tasksToDelete)
	*syncResponse = append(*syncResponse, SyncResponseItem{
		Type: "item_delete",
		Uuid: uuid.New().String(),
		Args: &SyncResponseArgs{
			Ids: idsToDelete,
		},
	})
}

func prepareDeferredSync(actions *[]adapters.TaskAction, syncResponse *[]SyncResponseItem) {
	tasksToDefer := getAllTasksWithAction(actions, adapters.ActionDefer)
	for _, task := range *tasksToDefer {
		// TODO: make these tags configurable
		tags := updateLabels(task.Task.Tags, "someday_maybe", []string{"next"})
		item := SyncResponseItem{
			Type: "item_update",
			Uuid: uuid.New().String(),
			Args: &SyncResponseArgs{
				Id:     &task.Task.ID,
				Labels: &tags,
			},
		}
		*syncResponse = append(*syncResponse, item)
	}
}

func prepareRevalidateSync(actions *[]adapters.TaskAction, syncResponse *[]SyncResponseItem) {
	tasksToRevalidate := getAllTasksWithAction(actions, adapters.ActionRevalidate)
	comment := "Revalidated on " + time.Now().Format("2006-01-02")
	for _, task := range *tasksToRevalidate {
		*syncResponse = append(*syncResponse, SyncResponseItem{
			Type: "note_add",
			Uuid: uuid.New().String(),
			Args: &SyncResponseArgs{
				ItemId:  &task.Task.ID,
				Content: &comment,
			},
		})
	}
}

func getAllTasksWithAction(actions *[]adapters.TaskAction, action adapters.Action) *[]adapters.TaskAction {
	var tasks []adapters.TaskAction
	for _, a := range *actions {
		if a.Action == action {
			tasks = append(tasks, a)
		}
	}
	return &tasks
}

func extractIds(tasks *[]adapters.TaskAction) *[]string {
	var ids []string
	for _, task := range *tasks {
		ids = append(ids, task.Task.ID)
	}
	return &ids
}

func updateLabels(labels []string, labelToAdd string, labelsToRemove []string) []string {
	labels = append(labels, labelToAdd)

	allKeys := make(map[string]bool)
	list := []string{}

	for _, item := range labelsToRemove {
		allKeys[item] = false
	}

	for _, item := range labels {
		if _, exists := allKeys[item]; !exists {
			allKeys[item] = true
		}
	}

	for label, toInclude := range allKeys {
		if toInclude {
			list = append(list, label)
		}
	}
	return list
}

func NewTodoistAdapter() (adapters.TaskManagerAdapter, error) {
	return &TodoistAdapter{}, nil
}
