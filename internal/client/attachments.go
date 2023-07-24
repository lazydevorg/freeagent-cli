package client

import "time"

type Attachment struct {
	Url              string    `json:"url"`
	ContentSrc       string    `json:"content_src"`
	ContentSrcMedium string    `json:"content_src_medium"`
	ContentSrcSmall  string    `json:"content_src_small"`
	ExpiresAt        time.Time `json:"expires_at"`
	ContentType      string    `json:"content_type"`
	FileName         string    `json:"file_name"`
	FileSize         int       `json:"file_size"`
	Description      string    `json:"description"`
}

func GetAttachment(id string) (*Attachment, error) {
	attachment, err := GetEntity[Attachment]("attachments/"+id, "attachment")
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return attachment, nil
}
