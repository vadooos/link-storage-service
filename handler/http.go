package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"link-storage-service/dto"
	"link-storage-service/repository"
	"link-storage-service/service"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Handler struct {
	svc *service.LinkService
}

func New(svc *service.LinkService) *Handler { return &Handler{svc: svc} }

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET     /links", h.links)
	mux.HandleFunc("POST    /links", h.createLink)
	mux.HandleFunc("GET     /links/{short_code}", h.getLinkByCode)
	mux.HandleFunc("DELETE  /links/{short_code}", h.deleteLink)
	mux.HandleFunc("GET     /links/{short_code}/stats", h.stats)
	return mux
}

func (h *Handler) links(w http.ResponseWriter, r *http.Request) {
	limit := 20
	offset := 0
	{
		n, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err == nil && n > 0 {
			limit = n
		}
	}
	{
		n, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err == nil && n >= 0 {
			offset = n
		}
	}

	links, err := h.svc.List(limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := dto.ListLinksResponse{Items: dto.LinkItemsFromModel(links), Limit: limit, Offset: offset}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) createLink(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateLinkRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "decode JSON error: "+err.Error())
		return
	}
	if err := validateURL(req.URL); err != nil {
		writeError(w, http.StatusBadRequest, "validate URL error: "+err.Error())
		return
	}

	code, err := h.svc.CreateLink(req.URL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "create link error: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, dto.CreateLinkResponse{ShortCode: code})
}

func (h *Handler) getLinkByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("short_code")
	link, err := h.svc.GetLinkAndIncrementVisits(code)
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, dto.GetLinkResponse{URL: link.OriginalURL, Visits: link.Visits})
}

func (h *Handler) deleteLink(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("short_code")
	err := h.svc.DeleteLink(code)
	if err != nil {
		respondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) stats(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("short_code")
	link, err := h.svc.GetLinkByShortCode(code)
	if err != nil {
		respondError(w, err)
		return
	}
	resp := dto.LinkStatsResponse{
		ShortCode: link.ShortCode,
		URL:       link.OriginalURL,
		Visits:    link.Visits,
		CreatedAt: link.CreatedAt}
	writeJSON(w, http.StatusOK, resp)
}

func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		return err
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Printf("encode JSON error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, dto.ErrorResponse{Error: msg})
}

func respondError(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeError(w, http.StatusInternalServerError, err.Error())
}

func validateURL(raw string) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fmt.Errorf("url is required")
	}
	u, err := url.ParseRequestURI(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("invalid url")
	}
	return nil
}
