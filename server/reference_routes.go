package server

//func (s Server) addReference(w http.ResponseWriter, r *http.Request) {
//	l := s.logger.With(slog.String("path", "POST /requirement"))
//	l.LogAttrs(r.Context(), slog.LevelInfo, "Handling create reference request")
//	reqBody := types.Reference{}
//	err := json.NewDecoder(r.Body).Decode(&reqBody)
//	if err != nil {
//		l.LogAttrs(r.Context(), slog.LevelWarn, "Error when decoding request from client", slog.String("error", err.Error()))
//		w.WriteHeader(http.StatusBadRequest)
//		r.Body.Close()
//		return
//	}
//	defer r.Body.Close()
//	requirement, err := s.backend.AddReference(reqBody)
//	if errors.Is(err, backend.ErrMissingArgs) {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	l.LogAttrs(r.Context(), slog.LevelInfo, "Encoding response back to client")
//	w.WriteHeader(http.StatusCreated)
//	w.Header().Set("Content-Type", "application/json")
//	err = json.NewEncoder(w).Encode(requirement)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		l.LogAttrs(r.Context(), slog.LevelError, "Error serializing reference to client", slog.String("error", err.Error()))
//		return
//	}
//}
//
//func (s Server) getReference(w http.ResponseWriter, r *http.Request) {
//	l := s.logger.With(slog.String("path", "GET /reference"))
//	l.LogAttrs(r.Context(), slog.LevelInfo, "Handling get reference request")
//	id := r.PathValue("id")
//	requirement, err := s.backend.GetReference(id)
//	if errors.Is(err, backend.ErrReferenceNotFound) {
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//	}
//	w.Header().Set("Content-Type", "application/json")
//	err = json.NewEncoder(w).Encode(requirement)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		l.LogAttrs(r.Context(), slog.LevelError, "Error when serializing reference to client", slog.String("error", err.Error()))
//		return
//	}
//}
//
//func (s Server) getReferences(w http.ResponseWriter, r *http.Request) {
//	l := s.logger.With(slog.String("path", "GET /references"))
//	l.LogAttrs(r.Context(), slog.LevelInfo, "Handling get references request")
//	refs, err := s.backend.GetReferences()
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	if len(refs) == 0 {
//		refs = make([]types.Reference, 0)
//	}
//	w.Header().Set("Content-Type", "application/json")
//	err = json.NewEncoder(w).Encode(refs)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		l.LogAttrs(r.Context(), slog.LevelError, "Error when serializing references to client", slog.String("error", err.Error()))
//		return
//	}
//}
