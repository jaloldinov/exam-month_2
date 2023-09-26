package handler

import (
	"market/pkg/logger"
	"market/storage"
)

type Handler struct {
	strg storage.StorageI
	log  logger.LoggerI
}

func NewHandler(strg storage.StorageI, loger logger.LoggerI) *Handler {
	return &Handler{strg: strg, log: loger}
}
