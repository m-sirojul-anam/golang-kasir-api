package handlers

import (
	"encoding/json"
	"kasir-api/dto"
	"kasir-api/models"
	"kasir-api/services"
	"kasir-api/utils"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/swaggo/swag/example/celler/httputil"
	_ "github.com/swaggo/swag/example/celler/model"
)

const (
	productPathPrefix       = "/api/product/"
	invalidProductIDMsg     = "Invalid product ID"
	contentTypeHeader       = "Content-Type"
	contentTypeJSON         = "application/json"
	failedEncodeResponseMsg = "Failed to encode response"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{productService: service}
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Get list of products
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /api/product [get]
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	products, err := h.productService.GetAllProducts(name)
	if err != nil {
		utils.WriteJSON(
			w, http.StatusInternalServerError, dto.ErrorResponse("failed to fetch products"),
		)
		return
	}

	if products == nil {
		products = []models.Product{}
	}

	utils.WriteJSON(w, http.StatusOK, dto.SuccessResponse("Success get products", products))
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product models.CreateProductRequest
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.productService.CreateProduct(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(contentTypeHeader, contentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, failedEncodeResponseMsg, http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r)
	case http.MethodPut:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, productPathPrefix)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, invalidProductIDMsg, http.StatusBadRequest)
		return
	}

	product, err := h.productService.GetProductByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set(contentTypeHeader, contentTypeJSON)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, failedEncodeResponseMsg, http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, productPathPrefix)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, invalidProductIDMsg, http.StatusBadRequest)
		return
	}

	var product models.CreateProductRequest
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product.ID = id
	err = h.productService.UpdateProduct(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(contentTypeHeader, contentTypeJSON)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, failedEncodeResponseMsg, http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, productPathPrefix)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, invalidProductIDMsg, http.StatusBadRequest)
		return
	}

	err = h.productService.DeleteProduct(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentTypeHeader, contentTypeJSON)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Product deleted successfully",
	}); err != nil {
		http.Error(w, failedEncodeResponseMsg, http.StatusInternalServerError)
		return
	}
}
