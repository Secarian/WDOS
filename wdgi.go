package main

import (
	"net/http"

	agi "imuslab.com/arozos/mod/agi"
	prout "imuslab.com/arozos/mod/prouter"
	"imuslab.com/arozos/mod/utils"
)

var (
	WDGIGateway *wdgi.Gateway
)

func WDGIInit() {
	//Create new WDGI Gateway object
	gw, err := wdgi.NewGateway(wdgi.WdgiSysInfo{
		BuildVersion:         build_version,
		InternalVersion:      internal_version,
		LoadedModule:         moduleHandler.GetModuleNameList(),
		ReservedTables:       []string{"auth", "permisson", "register", "desktop"},
		ModuleRegisterParser: moduleHandler.RegisterModuleFromAGI,
		PackageManager:       packageManager,
		UserHandler:          userHandler,
		StartupRoot:          "./app",
		ActivateScope:        []string{"./app", "./subservice"},
		FileSystemRender:     thumbRenderHandler,
		ShareManager:         shareManager,
		NightlyManager:       nightlyManager,
		TempFolderPath:       *tmp_directory,
	})
	if err != nil {
		systemWideLogger.PrintAndLog("WDGI", "WDGI Gateway Initialization Failed", err)
	}

	//Register user request handler endpoint
	http.HandleFunc("/system/wdgi/interface", func(w http.ResponseWriter, r *http.Request) {
		//Require login check
		authAgent.HandleCheckAuth(w, r, func(w http.ResponseWriter, r *http.Request) {
			//API Call from actual human users
			thisuser, _ := gw.Option.UserHandler.GetUserInfoFromRequest(w, r)
			gw.InterfaceHandler(w, r, thisuser)
		})
	})

	//Register external API request handler endpoint
	http.HandleFunc("/api/wdgi/interface", func(w http.ResponseWriter, r *http.Request) {
		//Check if token exists
		token, err := utils.PostPara(r, "token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Unauthorized (token is empty)"))
			return
		}

		//Validate Token
		if authAgent.TokenValid(token) {
			//Valid
			thisUsername, err := gw.Option.UserHandler.GetAuthAgent().GetTokenOwner(token)
			if err != nil {
				systemWideLogger.PrintAndLog("WDGI", "Unable to validate token owner", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Internal Server Error"))
				return
			}
			thisuser, _ := gw.Option.UserHandler.GetUserInfoFromUsername(thisUsername)
			gw.APIHandler(w, r, thisuser)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Unauthorized (Invalid / expired token)"))
			return
		}

	})

	http.HandleFunc("/api/wdgi/exec", gw.HandleAgiExecutionRequestWithToken)

	// external WDGI related function
	externalWDGIRouter := prout.NewModuleRouter(prout.RouterOption{
		ModuleName:  "WDE Serverless",
		AdminOnly:   false,
		UserHandler: userHandler,
		DeniedHandler: func(w http.ResponseWriter, r *http.Request) {
			errorHandlePermissionDenied(w, r)
		},
	})
	externalWDGIRouter.HandleFunc("/api/wdgi/listExt", gw.ListExternalEndpoint)
	externalWDGIRouter.HandleFunc("/api/wdgi/addExt", gw.AddExternalEndPoint)
	externalWDGIRouter.HandleFunc("/api/wdgi/rmExt", gw.RemoveExternalEndPoint)

	WDGIGateway = gw
}
