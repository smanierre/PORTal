package api

import "net/http"

func (s Server) addMemberQualification(w http.ResponseWriter, r *http.Request) {
	memberID := r.PathValue("id")
	qualID := r.PathValue("qualID")
	err := s.backend.AddMemberQualification(qualID, memberID)
	if err != nil {
		panic(err)
	}
}

func (s Server) deleteMemberQualification(w http.ResponseWriter, r *http.Request) {
	memberID := r.PathValue("id")
	qualID := r.PathValue("qualID")
	err := s.backend.DeleteMemberQualification(qualID, memberID)
	if err != nil {
		panic(err)
	}
}
