package todoist

import (
	"github.com/dormunis/gitd/adapters"
	"time"
)

func (t *TodoistSyncResponse) ToTasks() []adapters.Task {
	var tasks []adapters.Task
	for _, item := range *t.Items {
		updatedDate := getLastNoteDateFromItem(*item.ID, t.Notes)
		if updatedDate != nil && updatedDate.After(*item.AddedAt) {
			updatedDate = item.UpdatedAt
		} else {
			updatedDate = item.AddedAt
		}

		tasks = append(tasks, adapters.Task{
			ID:          *item.ID,
			Project:     t.getProjectName(*item.ProjectID),
			Content:     *item.Content,
			CreatedDate: *item.AddedAt,
			UpdatedDate: *updatedDate,
			Tags:        *item.Labels,
			TaskManger:  "todoist",
			Status:      deriveStatus(item),
			Priority:    adapters.Priority(*item.Priority),
		})
	}
	return tasks
}

func (t *TodoistSyncResponse) getProjectName(projectID string) string {
	for _, project := range *t.Projects {
		if *project.ID == projectID {
			return *project.Name
		}
	}
	return "<Unknown>"
}

func getLastNoteDateFromItem(itemID string, notes *[]Note) *time.Time {
	var latest *time.Time
	for _, note := range *notes {
		if *note.ItemID == itemID && (latest == nil || note.PostedAt.After(*latest)) {
			latest = note.PostedAt
		}
	}
	if latest != nil {
		return latest
	}
	return nil
}

func deriveStatus(item Item) adapters.Status {
	if item.CompletedAt != nil {
		return adapters.StatusCompleted
	}
	// TODO: make these labels configurable in TodoistConfig or something
	for _, label := range *item.Labels {
		if label == "next" {
			return adapters.StatusNext
		}
		if label == "someday_maybe" {
			return adapters.StatusSomeday
		}
	}
	return adapters.StatusActive
}

type TodoistSyncRequest struct {
	SyncToken     string   `json:"sync_token"`
	ResourceTypes []string `json:"resource_types"`
}

type TodoistSyncResponse struct {
	Items    *[]Item    `json:"items"`
	Labels   *[]Label   `json:"labels"`
	Notes    *[]Note    `json:"notes"`
	Projects *[]Project `json:"projects"`
	Sections *[]Section `json:"sections"`
	User     *User      `json:"user"`
}

type Item struct {
	AddedAt     *time.Time `json:"added_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Content     *string    `json:"content"`
	Description *string    `json:"description"`
	ID          *string    `json:"id"`
	Labels      *[]string  `json:"labels"`
	ParentID    *string    `json:"parent_id"`
	Priority    *int       `json:"priority"`
	ProjectID   *string    `json:"project_id"`
	SectionID   *string    `json:"section_id"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type Label struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
}

type Note struct {
	Content        *string         `json:"content"`
	FileAttachment *FileAttachment `json:"file_attachment"`
	ID             *string         `json:"id"`
	ItemID         *string         `json:"item_id"`
	PostedAt       *time.Time      `json:"posted_at"`
}

type Project struct {
	CreatedAt *time.Time `json:"created_at"`
	ID        *string    `json:"id"`
	Name      *string    `json:"name"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type Section struct {
	AddedAt   *time.Time `json:"added_at"`
	ID        *string    `json:"id"`
	Name      *string    `json:"name"`
	ProjectID *string    `json:"project_id"`
}

type FileAttachment struct {
	ResourceType *string `json:"resource_type"`
	Title        *string `json:"title"`
	URL          *string `json:"url"`
}

type User struct {
	AvatarBig      *string `json:"avatar_big"`
	AvatarMedium   *string `json:"avatar_medium"`
	AvatarS640     *string `json:"avatar_s640"`
	AvatarSmall    *string `json:"avatar_small"`
	Email          *string `json:"email"`
	FullName       *string `json:"full_name"`
	ID             *string `json:"id"`
	ImageID        *string `json:"image_id"`
	InboxProjectID *string `json:"inbox_project_id"`
	TzInfo         *TzInfo `json:"tz_info"`
}

type TzInfo struct {
	GmtString *string `json:"gmt_string"`
	Hours     *int    `json:"hours"`
	IsDst     *int    `json:"is_dst"`
	Minutes   *int    `json:"minutes"`
	Timezone  *string `json:"timezone"`
}

type Due struct {
	Date        *string    `json:"date"`
	IsRecurring *bool      `json:"is_recurring"`
	Datetime    *time.Time `json:"datetime"`
	String      *string    `json:"string"`
	Timezone    *string    `json:"timezone"`
}

type Duration struct {
	Amount *int    `json:"amount"`
	Unit   *string `json:"unit"`
}

type SyncResponseItem struct {
	Type   string            `json:"type"`
	Uuid   string            `json:"uuid"`
	TempId *string           `json:"temp_id,omitempty"` // used if the item is not yet created
	Args   *SyncResponseArgs `json:"args"`
}

type SyncResponseArgs struct {
	Id      *string   `json:"id,omitempty"`
	ItemId  *string   `json:"item_id,omitempty"`
	Ids     *[]string `json:"ids,omitempty"`
	Content *string   `json:"content,omitempty"`
	Labels  *[]string `json:"labels,omitempty"`
}
