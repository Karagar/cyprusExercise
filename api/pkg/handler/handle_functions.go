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

// getHandlerFunc is for selection the appropriate handler according to the config
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

// getCompanyHandler is for reading company rows
func getCompanyHandler(h *Handler, w http.ResponseWriter, r *http.Request) {
	query := getRowsByParams(h, r)
	params := r.URL.Query()

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

// postCompanyHandler is for creating new company rows with client object
func postCompanyHandler(h *Handler, w http.ResponseWriter, r *http.Request) {
	company := &([]structs.Company{})
	query := h.DB.WithContext(h.ctx).Table("company")

	err := utils.ReadJsonBody(r.Body, company)
	if err != nil {
		h.handleProblems(w, http.StatusBadRequest, err)
		return
	}

	result := query.Create(&company)
	if err != nil {
		h.handleProblems(w, http.StatusBadGateway, result.Error)
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

// putCompanyHandler is for updating company rows with client object
func putCompanyHandler(h *Handler, w http.ResponseWriter, r *http.Request) {
	query := getRowsByParams(h, r)
	company := &(structs.Company{})
	response := structs.CompanyResponse{}

	// select companies
	query = query.Find(&response.Data)
	if query.Error != nil {
		h.handleProblems(w, http.StatusInternalServerError, query.Error)
		return
	}

	// get selected companies uuids
	uuidList := [][]byte{}
	for _, v := range response.Data {
		uuidList = append(uuidList, v.ID)
	}

	// get client's put object
	err := utils.ReadJsonBody(r.Body, company)
	if err != nil {
		h.handleProblems(w, http.StatusBadRequest, err)
		return
	}

	// make update on all fields
	filterValues := getFilterList()
	fieldColumnMap := getFieldColumnMap()
	fields := make([]string, 0, len(filterValues))
	for _, v := range filterValues {
		fields = append(fields, fieldColumnMap[v])
	}
	fields = append(fields, fieldColumnMap["DTUpdated"])
	dt_updated := time.Now().UTC().Format(time.RFC3339)
	company.DTUpdated = &dt_updated
	query.Select(fields).Updates(company)

	// get updated companies by uuids
	query = h.DB.WithContext(h.ctx).Table("company")
	query = query.Where(uuidList).Find(&response.Data)
	sendData(h, w, r, query, response)
}

// patchCompanyHandler is for merge company rows with client object
func patchCompanyHandler(h *Handler, w http.ResponseWriter, r *http.Request) {
	query := getRowsByParams(h, r)
	company := &(structs.Company{})
	response := structs.CompanyResponse{}

	// select companies
	query = query.Find(&response.Data)
	if query.Error != nil {
		h.handleProblems(w, http.StatusInternalServerError, query.Error)
		return
	}

	// get selected companies uuids
	uuidList := [][]byte{}
	for _, v := range response.Data {
		uuidList = append(uuidList, v.ID)
	}

	// make patch
	err := utils.ReadJsonBody(r.Body, company)
	if err != nil {
		h.handleProblems(w, http.StatusBadRequest, err)
		return
	}
	dt_updated := time.Now().UTC().Format(time.RFC3339)
	company.DTUpdated = &dt_updated
	query.Updates(company)

	// get updated companies by uuids
	query = h.DB.WithContext(h.ctx).Table("company")
	query = query.Where(uuidList).Find(&response.Data, uuidList)
	sendData(h, w, r, query, response)
}

// deleteCompanyHandler is for archivating company rows in DB
func deleteCompanyHandler(h *Handler, w http.ResponseWriter, r *http.Request) {
	query := getRowsByParams(h, r)
	archive := true
	dt_archive := time.Now().UTC().Format(time.RFC3339)
	response := structs.CompanyResponse{}
	query = query.Find(&response.Data)

	query.Updates(&structs.Company{
		Archive:    &archive,
		DTArchived: &dt_archive,
		DTUpdated:  &dt_archive,
	})

	sendData(h, w, r, query, response)
}

// getRowsByParams filtered query in accordance with client request
func getRowsByParams(h *Handler, r *http.Request) *gorm.DB {
	filterValues := getFilterList()
	fieldColumnMap := getFieldColumnMap()
	company := map[string]interface{}{
		fieldColumnMap["Archive"]: false,
	}

	query := h.DB.WithContext(h.ctx).Table("company")

	params := r.URL.Query()
	for k, _ := range params {
		if slices.Contains(filterValues, k) {
			company[fieldColumnMap[k]] = params.Get(k)
		}
	}

	identifier := params.Get("Uuid")
	if identifier != "" {
		_, err := uuid.Parse(identifier)
		if err == nil {
			company[fieldColumnMap["Uuid"]] = identifier
		} else {
			company[fieldColumnMap["Uuid"]] = defaultIdentidier
		}
	}
	return query.Where(company)
}

// sendData returns the result to the client
func sendData(h *Handler, w http.ResponseWriter, r *http.Request, query *gorm.DB, response structs.CompanyResponse) {
	if query.Error != nil {
		h.handleProblems(w, http.StatusInternalServerError, query.Error)
		return
	}

	for _, v := range response.Data {
		v.Uuid = utils.HandleUuid(v.ID)
	}
	response.Count = int(query.RowsAffected)
	body, err := json.Marshal(response)
	if err != nil {
		h.handleProblems(w, http.StatusInternalServerError, err)
		return
	}
	http.ServeContent(w, r, "index.json", time.Time{}, bytes.NewReader(body))
}

// getFieldColumnMap return map of Company fields - DB columns
func getFieldColumnMap() map[string]string {
	return map[string]string{
		"Uuid":      "id",
		"Name":      "company_name",
		"Code":      "code",
		"Country":   "country",
		"Website":   "website",
		"Phone":     "phone",
		"Archive":   "archive",
		"DTUpdated": "dt_updated",
	}
}

// getFilterList return list of possible client filters
func getFilterList() []string {
	return []string{"Name", "Code", "Country", "Website", "Phone"}
}
