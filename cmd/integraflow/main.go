// IntegraFlow IMS
//
// @title       IntegraFlow IMS API
// @version     1.0
// @description Integrated Management System API for risks, incidents, audits and CAPA.
//
// @contact.name  IntegraFlow IMS Team
// @contact.email ims@example.com
//
// @BasePath /
package main

import (
	"log"
	"net/http"
	"os"

	repoSqlite "github.com/xenakil/integraflow-ims/internal/repository/sqlite"
	"github.com/xenakil/integraflow-ims/internal/service"
	"github.com/xenakil/integraflow-ims/internal/transport/httpapi"
)

func main() {
	// Initialize SQLite DB
	db, err := repoSqlite.NewDB("integraflow.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	// Initialize repositories backed by SQLite
	riskRepo := repoSqlite.NewRiskRepository(db)
	incidentRepo := repoSqlite.NewIncidentRepository(db)
	auditRepo := repoSqlite.NewAuditRepository(db)
	actionRepo := repoSqlite.NewActionRepository(db)

	// Initialize services
	riskSvc := service.NewRiskService(riskRepo)
	incidentSvc := service.NewIncidentService(incidentRepo, riskRepo)
	auditSvc := service.NewAuditService(auditRepo)
	actionSvc := service.NewActionService(actionRepo, riskRepo, incidentRepo, auditRepo)
	dashboardSvc := service.NewDashboardService(riskRepo, incidentRepo, actionRepo)

	// HTTP API server
	server := httpapi.NewServer(riskSvc, incidentSvc, auditSvc, actionSvc, dashboardSvc)

	port := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}

	log.Printf("IntegraFlow IMS API (SQLite) listening on http://localhost%v\n", port)
	if err := http.ListenAndServe(port, server); err != nil {
		log.Fatal(err)
	}
}
