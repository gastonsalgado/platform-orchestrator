package controllers

import (
	"github.com/gastonsalgado/platform-orchestrator/backend/internal/managers"
	"go.uber.org/zap"
)

var Logger *zap.Logger
var gitManager = managers.GetGitManagerInstance()
