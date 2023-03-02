package route

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/swaggo/http-swagger"

	delivery "github.com/openinfradev/tks-api/internal/delivery/http"
	"github.com/openinfradev/tks-api/internal/helper"
	"github.com/openinfradev/tks-api/internal/repository"
	"github.com/openinfradev/tks-api/internal/usecase"
	argowf "github.com/openinfradev/tks-api/pkg/argo-client"
	"github.com/openinfradev/tks-api/pkg/log"
)

const (
	API_VERSION = "/1.0"
	API_PREFIX  = "/api"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func SetupRouter(db *gorm.DB, argoClient argowf.ArgoClient, asset http.Handler) http.Handler {
	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	// [TODO] Transaction
	//r.Use(transactionMiddleware(db))

	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	/*
		// 1.0
		r.HandleFunc("/api/1.0/signin", handler.Signin).Methods(http.MethodPost)
		r.HandleFunc("/api/1.0/signout", handler.Signout).Methods(http.MethodPost)
		r.HandleFunc("/api/1.0/token", handler.SigninByToken).Methods(http.MethodPost)
		r.Handle("/api/1.0/signup", authMiddleware(http.HandlerFunc(handler.Signup))).Methods(http.MethodPost)

		r.Handle("/api/1.0/users", authMiddleware(http.HandlerFunc(handler.GetUsers))).Methods(http.MethodGet)
		r.Handle("/api/1.0/users/{userId}/password", authMiddleware(http.HandlerFunc(handler.UpdatePassword))).Methods(http.MethodPut)
		r.Handle("/api/1.0/users/{userId}/role", authMiddleware(http.HandlerFunc(handler.UpdateRole))).Methods(http.MethodPut)

		r.Handle("/api/1.0/settings/server", authMiddleware(http.HandlerFunc(handler.GetServerSettings))).Methods(http.MethodGet)

		r.Handle("/api/1.0/dashboard", authMiddleware(http.HandlerFunc(handler.GetOverview))).Methods(http.MethodGet)
		r.Handle("/api/1.0/dashboard/kube-info", authMiddleware(http.HandlerFunc(handler.GetKubernetesInfo))).Methods(http.MethodGet)
		r.Handle("/api/1.0/dashboard/kube-events", authMiddleware(http.HandlerFunc(handler.GetAdnormalKubernetesEvents))).Methods(http.MethodGet)
		r.Handle("/api/1.0/dashboard/kube-pods", authMiddleware(http.HandlerFunc(handler.GetAdnormalKubernetesPods))).Methods(http.MethodGet)

		r.Handle("/api/1.0/organizations", authMiddleware(http.HandlerFunc(handler.CreateOrganization))).Methods(http.MethodPost)
		r.Handle("/api/1.0/organizations", authMiddleware(http.HandlerFunc(handler.GetOrganizations))).Methods(http.MethodGet)
		r.Handle("/api/1.0/organizations/users", authMiddleware(http.HandlerFunc(handler.GetOrganizationUsers))).Methods(http.MethodGet)
		r.Handle("/api/1.0/organizations/{organizationId}", authMiddleware(http.HandlerFunc(handler.GetOrganization))).Methods(http.MethodGet)
		r.Handle("/api/1.0/organizations/{organizationId}/user", authMiddleware(http.HandlerFunc(handler.AddOrganizationUser))).Methods(http.MethodPost)
		r.Handle("/api/1.0/organizations/{organizationId}/user", authMiddleware(http.HandlerFunc(handler.RemoveOrganizationUser))).Methods(http.MethodDelete)

		r.Handle("/api/1.0/clusters", authMiddleware(http.HandlerFunc(handler.GetClusters))).Methods(http.MethodGet)
		r.Handle("/api/1.0/clusters", authMiddleware(http.HandlerFunc(handler.CreateCluster))).Methods(http.MethodPost)
		r.Handle("/api/1.0/clusters/{clusterId}", authMiddleware(http.HandlerFunc(handler.GetCluster))).Methods(http.MethodGet)
		r.Handle("/api/1.0/clusters/{clusterId}/kubeconfig", authMiddleware(http.HandlerFunc(handler.GetClusterKubeConfig))).Methods(http.MethodGet)
		r.Handle("/api/1.0/clusters/{clusterId}/kube-info", authMiddleware(http.HandlerFunc(handler.GetKubernetesInfo))).Methods(http.MethodGet)
		r.Handle("/api/1.0/clusters/{clusterId}/applications", authMiddleware(http.HandlerFunc(handler.GetClusterApplications))).Methods(http.MethodGet)
		r.Handle("/api/1.0/clusters/{clusterId}/applications/kube-info", authMiddleware(http.HandlerFunc(handler.GetClusterApplicationsKubeInfo))).Methods(http.MethodGet)
		r.Handle("/api/1.0/clusters/{clusterId}", authMiddleware(http.HandlerFunc(handler.DeleteCluster))).Methods(http.MethodDelete)
		r.Handle("/api/1.0/clusters/{clusterId}/kube-resources", authMiddleware(http.HandlerFunc(handler.GetClusterKubeResources))).Methods(http.MethodGet)
		r.Handle("/api/1.0/clusters/{clusterId}/kube-resources/{namespace}/istio-label", authMiddleware(http.HandlerFunc(handler.SetIstioLabel))).Methods(http.MethodPost)

		r.Handle("/api/1.0/infra-providers", authMiddleware(http.HandlerFunc(handler.GetInfraProviders))).Methods(http.MethodGet)
		r.Handle("/api/1.0/infra-providers", authMiddleware(http.HandlerFunc(handler.CreateInfraProvider))).Methods(http.MethodPost)
		r.Handle("/api/1.0/infra-providers/{infraProviderId}", authMiddleware(http.HandlerFunc(handler.GetInfraProvider))).Methods(http.MethodGet)

		r.Handle("/api/1.0/applications", authMiddleware(http.HandlerFunc(handler.GetApplications))).Methods(http.MethodGet)
		r.Handle("/api/1.0/applications", authMiddleware(http.HandlerFunc(handler.CreateApplication))).Methods(http.MethodPost)
		r.Handle("/api/1.0/applications/{applicationId}", authMiddleware(http.HandlerFunc(handler.GetApplication))).Methods(http.MethodGet)
		r.Handle("/api/1.0/applications/{applicationId}", authMiddleware(http.HandlerFunc(handler.DeleteApplication))).Methods(http.MethodDelete)
		r.Handle("/api/1.0/applications/{applicationId}/kube-info", authMiddleware(http.HandlerFunc(handler.GetApplicationKubeInfo))).Methods(http.MethodGet)

		r.Handle("/api/1.0/app-serve-apps", authMiddleware(http.HandlerFunc(handler.GetAppServeApps))).Methods(http.MethodGet)
		r.Handle("/api/1.0/app-serve-apps", authMiddleware(http.HandlerFunc(handler.CreateAppServeApp))).Methods(http.MethodPost)
		r.Handle("/api/1.0/app-serve-apps/{asaId}", authMiddleware(http.HandlerFunc(handler.UpdateAppServeApp))).Methods(http.MethodPut)
		r.Handle("/api/1.0/app-serve-apps/{asaId}", authMiddleware(http.HandlerFunc(handler.GetAppServeApp))).Methods(http.MethodGet)
		r.Handle("/api/1.0/app-serve-apps/{asaId}", authMiddleware(http.HandlerFunc(handler.DeleteAppServeApp))).Methods(http.MethodDelete)


		r.Handle("/api/1.0/stacks", authMiddleware(http.HandlerFunc(handler.GetStacks))).Methods(http.MethodGet)
		r.Handle("/api/1.0/stacks", authMiddleware(http.HandlerFunc(handler.CreateStack))).Methods(http.MethodPost)
		r.Handle("/api/1.0/stacks/{clusterId}", authMiddleware(http.HandlerFunc(handler.GetStack))).Methods(http.MethodGet)
		r.Handle("/api/1.0/stacks/{clusterId}", authMiddleware(http.HandlerFunc(handler.DeleteStack))).Methods(http.MethodDelete)
	*/

	authHandler := delivery.NewAuthHandler(usecase.NewAuthUsecase(repository.NewAuthRepository(db)))
	r.HandleFunc(API_PREFIX+API_VERSION+"/auth/signin", authHandler.Signin).Methods(http.MethodPost)
	r.HandleFunc(API_PREFIX+API_VERSION+"/auth/signup", authHandler.Signup).Methods(http.MethodPost)

	organizationHandler := delivery.NewOrganizationHandler(usecase.NewOrganizationUsecase(repository.NewOrganizationRepository(db), argoClient))
	r.Handle(API_PREFIX+API_VERSION+"/organizations", authMiddleware(http.HandlerFunc(organizationHandler.CreateOrganization))).Methods(http.MethodPost)
	r.Handle(API_PREFIX+API_VERSION+"/organizations", authMiddleware(http.HandlerFunc(organizationHandler.GetOrganizations))).Methods(http.MethodGet)
	r.Handle(API_PREFIX+API_VERSION+"/organizations/{organizationId}", authMiddleware(http.HandlerFunc(organizationHandler.GetOrganization))).Methods(http.MethodGet)
	r.Handle(API_PREFIX+API_VERSION+"/organizations/{organizationId}", authMiddleware(http.HandlerFunc(organizationHandler.DeleteOrganization))).Methods(http.MethodDelete)

	clusterHandler := delivery.NewClusterHandler(usecase.NewClusterUsecase(
		repository.NewClusterRepository(db),
		repository.NewAppGroupRepository(db),
		argoClient))
	r.Handle(API_PREFIX+API_VERSION+"/clusters", authMiddleware(http.HandlerFunc(clusterHandler.CreateCluster))).Methods(http.MethodPost)
	r.Handle(API_PREFIX+API_VERSION+"/clusters", authMiddleware(http.HandlerFunc(clusterHandler.GetClusters))).Methods(http.MethodGet)
	r.Handle(API_PREFIX+API_VERSION+"/clusters/{clusterId}", authMiddleware(http.HandlerFunc(clusterHandler.GetCluster))).Methods(http.MethodGet)
	r.Handle(API_PREFIX+API_VERSION+"/clusters/{clusterId}", authMiddleware(http.HandlerFunc(clusterHandler.DeleteCluster))).Methods(http.MethodDelete)

	appGroupHandler := delivery.NewAppGroupHandler(usecase.NewAppGroupUsecase(
		repository.NewAppGroupRepository(db),
		repository.NewClusterRepository(db),
		argoClient))
	r.Handle(API_PREFIX+API_VERSION+"/app-groups", authMiddleware(http.HandlerFunc(appGroupHandler.CreateAppGroup))).Methods(http.MethodPost)
	r.Handle(API_PREFIX+API_VERSION+"/app-groups", authMiddleware(http.HandlerFunc(appGroupHandler.GetAppGroups))).Methods(http.MethodGet)
	r.Handle(API_PREFIX+API_VERSION+"/app-groups/{appGroupId}", authMiddleware(http.HandlerFunc(appGroupHandler.GetAppGroup))).Methods(http.MethodGet)
	r.Handle(API_PREFIX+API_VERSION+"/app-groups/{appGroupId}", authMiddleware(http.HandlerFunc(appGroupHandler.DeleteAppGroup))).Methods(http.MethodDelete)

	appServeAppHandler := delivery.NewAppServeAppHandler(usecase.NewAppServeAppUsecase(
		repository.NewAppServeAppRepository(db),
		argoClient))
	r.Handle(API_PREFIX+API_VERSION+"/app-serve-apps", authMiddleware(http.HandlerFunc(appServeAppHandler.CreateAppServeApp))).Methods(http.MethodPost)
	r.Handle(API_PREFIX+API_VERSION+"/app-serve-apps", authMiddleware(http.HandlerFunc(appServeAppHandler.GetAppServeApps))).Methods(http.MethodGet)
	r.Handle(API_PREFIX+API_VERSION+"/app-serve-apps/{appServeAppId}", authMiddleware(http.HandlerFunc(appServeAppHandler.GetAppServeApp))).Methods(http.MethodGet)
	r.Handle(API_PREFIX+API_VERSION+"/app-serve-apps/{appServeAppId}", authMiddleware(http.HandlerFunc(appServeAppHandler.DeleteAppServeApp))).Methods(http.MethodDelete)
	r.Handle(API_PREFIX+API_VERSION+"/app-serve-apps/{appServeAppId}", authMiddleware(http.HandlerFunc(appServeAppHandler.UpdateAppServeApp))).Methods(http.MethodPut)

	historyHandler := delivery.NewHistoryHandler(usecase.NewHistoryUsecase(repository.NewHistoryRepository(db)))
	r.Handle(API_PREFIX+API_VERSION+"/histories", authMiddleware(http.HandlerFunc(historyHandler.GetHistories))).Methods(http.MethodGet)

	// assets
	r.PathPrefix("/api/").HandlerFunc(http.NotFound)
	r.PathPrefix("/").Handler(asset).Methods(http.MethodGet)

	//withLog := handlers.LoggingHandler(os.Stdout, r)

	credentials := handlers.AllowCredentials()
	headersOk := handlers.AllowedHeaders([]string{"content-type", "Authorization", "Authorization-Type"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	return handlers.CORS(credentials, headersOk, originsOk, methodsOk)(r)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Possible values : "basic", "keycloak"
		authType := r.Header.Get("Authorization-Type")

		switch authType {
		case "keycloak":
			// [TODO] implementaion keycloak process
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Need implementation for keycloak"))
			return
		case "basic":
		default:
			tokenString := r.Header.Get("Authorization")
			if len(tokenString) == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Missing Authorization Header"))
				return
			}
			tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
			token, err := helper.VerifyToken(tokenString)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Error verifying JWT token: " + err.Error()))
				return
			}

			log.Info("[authMiddleware] accountId : ", token.Claims.(jwt.MapClaims)["AccountId"])
			log.Info("[authMiddleware] Id : ", token.Claims.(jwt.MapClaims)["Id"])
			accountId := token.Claims.(jwt.MapClaims)["AccountId"]
			id := token.Claims.(jwt.MapClaims)["ID"]
			r.Header.Set("AccountId", fmt.Sprint(accountId))
			r.Header.Set("ID", fmt.Sprint(id))
		}

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("***** START [%s %s] ***** ", r.Method, r.RequestURI))
		next.ServeHTTP(w, r)
		log.Info("***** END *****")
	})
}

func transactionMiddleware(db *gorm.DB) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			txHandle := db.Begin()
			log.Debug("beginning database transaction")

			defer func() {
				if r := recover(); r != nil {
					txHandle.Rollback()
				}
			}()

			recorder := &StatusRecorder{
				ResponseWriter: w,
				Status:         200,
			}

			r = r.WithContext(context.WithValue(ctx, "txHandle", txHandle))
			next.ServeHTTP(recorder, r)

			if StatusInList(recorder.Status, []int{http.StatusOK}) {
				log.Debug("committing transactions")
				if err := txHandle.Commit().Error; err != nil {
					log.Debug("trx commit error: ", err)
				}
			} else {
				log.Debug("rolling back transaction due to status code: ", recorder.Status)
				txHandle.Rollback()
			}
		})
	}
}

// StatusInList -> checks if the given status is in the list
func StatusInList(status int, statusList []int) bool {
	for _, i := range statusList {
		if i == status {
			return true
		}
	}
	return false
}
