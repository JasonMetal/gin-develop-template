package entity

type StsToken struct {
	SecretID     string   `json:"secret_id"`
	SecretKey    string   `json:"secret_key"`
	SessionToken string   `json:"session_token"`
	FileList     []string `json:"file_list"`
	Bucket       string   `json:"bucket"`
	Region       string   `json:"region"`
}
