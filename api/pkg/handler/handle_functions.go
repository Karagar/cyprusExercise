package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Karagar/cyprusExercise/pkg/structs"
	"github.com/Karagar/cyprusExercise/pkg/utils"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type HandlerFunc func(h *Handler, w http.ResponseWriter, r *http.Request)

const defaultIdentidier = "00000000-0000-0000-0000-000000000000"

func getHandlerFunc(funcName string) HandlerFunc {
	targetFunc := map[string]HandlerFunc{
		"getCompanyHandler":    getCompanyHandler,
		"postCompanyHandler":   postCompanyHandler,
		"putCompanyHandler":    putCompanyHandler,
		"patchCompanyHandler":  patchCompanyHandler,
		"deleteCompanyHandler": deleteCompanyHandler,
	}
	return targetFunc[funcName]
}

func getCompanyHandler(h *Handler, w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Serve ", h.Route.Path.URL, " (", h.Route.Path.Method, ")")
	filterValues := getFilterList()
	filterMap := getFilterMap()

	query := h.DB.WithContext(h.ctx).Table("company").Where("archive = ?", "False")

	params := r.URL.Query()
	for k, _ := range params {
		if slices.Contains(filterValues, k) {
			query = query.Where(filterMap[k]+" = ?", params.Get(k))
		}
	}

	identifier := params.Get("Uuid")
	if identifier != "" {
		_, err := uuid.Parse(identifier)
		if err == nil {
			query = query.Where(filterMap["Uuid"]+" = ?", identifier)
		} else {
			query = query.Where(filterMap["Uuid"]+" = ?", defaultIdentidier)
		}
	}

	limitStr := params.Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err == nil {
		query = query.Limit(limit)
	}

	offsetStr := params.Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err == nil {
		query = query.Offset(offset)
	}

	response := structs.CompanyResponse{}
	query = query.Find(&response.Data)
	sendData(h, w, r, query, response)
}

func postCompanyHandler(h *Handler, w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Serve ", h.Route.Path.URL, " (", h.Route.Path.Method, ")")
	company := &([]structs.Company{})
	query := h.DB.WithContext(h.ctx).Table("company")

	err := utils.ReadJsonBody(r.Body, company)
	if err != nil {
		h.handleProblems(w, err)
		return
	}

	result := query.Create(&company)
	if err != nil {
		h.handleProblems(w, result.Error)
		return
	}

	UUIDs := [][]byte{}
	for _, v := range *company {
		UUIDs = append(UUIDs, v.ID)
	}
	response := structs.CompanyResponse{}
	query = query.Find(&response.Data, UUIDs)
	sendData(h, w, r, query, response)
}

func putCompanyHandler(h *Handler, w http.ResponseWriter, r *http.Request) {}

func patchCompanyHandler(h *Handler, w http.ResponseWriter, r *http.Request) {}

func deleteCompanyHandler(h *Handler, w http.ResponseWriter, r *http.Request) {}

func getFilterMap() map[string]string {
	return map[string]string{
		"Uuid":    "id",
		"Name":    "company_name",
		"Code":    "code",
		"Country": "country",
		"Website": "website",
		"Phone":   "phone",
	}
}

func getFilterList() []string {
	return []string{"Name", "Code", "Country", "Website", "Phone"}
}

func sendData(h *Handler, w http.ResponseWriter, r *http.Request, query *gorm.DB, response structs.CompanyResponse) {
	if query.Error != nil {
		h.handleProblems(w, query.Error)
		return
	}

	for _, v := range response.Data {
		v.Uuid = utils.HandleUuid(v.ID)
	}
	response.Count = int(query.RowsAffected)
	body, err := json.Marshal(response)
	if err != nil {
		h.handleProblems(w, err)
		return
	}
	http.ServeContent(w, r, "index.json", time.Time{}, bytes.NewReader(body))
}
