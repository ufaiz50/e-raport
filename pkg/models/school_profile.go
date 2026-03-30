package models

type UpsertSchoolProfile struct {
	SchoolID       *string `json:"school_id"`
	SchoolName     string  `json:"school_name" binding:"required"`
	NPSN           string  `json:"npsn"`
	Address        string  `json:"address"`
	PrincipalName  string  `json:"principal_name"`
	PrincipalNIP   string  `json:"principal_nip"`
	HeadmasterSign string  `json:"headmaster_sign"`
	SchoolStamp    string  `json:"school_stamp"`
}
