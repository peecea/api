package api

import (
	"duval/internal/authentication"
	"duval/internal/route/docs"
	"duval/pkg/address"
	"duval/pkg/education"
	"duval/pkg/mark"
	"duval/pkg/media"
	cvtype "duval/pkg/media/cv"
	"duval/pkg/media/profile"
	"duval/pkg/media/video"
	"duval/pkg/planning"
	"duval/pkg/user"
	"net/http"
)

var Routes = []docs.RouteDocumentation{
	{
		HttpMethod:   http.MethodPost,
		RelativePath: "/upload",
		Handler:      media.Upload,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodHead,
		RelativePath: "/public",
		DocRoot:      "public",
		NeedToken:    false,
	},
	{
		HttpMethod:   http.MethodPost,
		RelativePath: "/register/:as",
		Handler:      user.Register,
		NeedToken:    false,
	},
	{
		HttpMethod:   http.MethodPost,
		RelativePath: "/register/:as/:email",
		Handler:      user.RegisterByEmail,
		NeedToken:    false,
	},
	{
		HttpMethod:   http.MethodPost,
		RelativePath: "/code/send",
		Handler:      user.SendUserEmailValidationCode,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodPost,
		RelativePath: "/code/verification/:code",
		Handler:      user.VerifyUserEmailValidationCode,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodGet,
		RelativePath: "/code",
		Handler:      user.GetCode,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodPost,
		RelativePath: "/login",
		Handler:      user.Login,
		NeedToken:    false,
	},
	{
		HttpMethod:   http.MethodPost,
		RelativePath: "/password",
		Handler:      user.NewPassword,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodGet,
		RelativePath: "/password/history",
		Handler:      user.GetPasswordHistory,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodGet,
		RelativePath: "/profile",
		Handler:      user.MyProfile,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodPut,
		RelativePath: "/profile/active",
		Handler:      user.ActivateUser,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodPut,
		RelativePath: "/profile",
		Handler:      user.UpdMyProfile,
		NeedToken:    true,
	},
	// Address route
	{
		HttpMethod:   http.MethodPost,
		RelativePath: "/address",
		Handler:      address.NewAddress,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodPut,
		RelativePath: "/address",
		Handler:      address.UpdateUserAddress,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodGet,
		RelativePath: "/address",
		Handler:      address.GetUserAddress,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodDelete,
		RelativePath: "/address",
		Handler:      address.RemoveUserAddress,
		NeedToken:    true,
	},
	//Profile image routes
	{
		HttpMethod:   http.MethodPost,
		RelativePath: "/profile/image",
		Handler:      profile.Upload,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodPut,
		RelativePath: "/profile/image",
		Handler:      profile.UpdateProfileImage,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodGet,
		RelativePath: "/profile/image",
		Handler:      profile.GetProfileImage,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodGet,
		RelativePath: "/profile/thumb",
		Handler:      profile.GetProfileThumb,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodDelete,
		RelativePath: "/profile/image",
		Handler:      profile.RemoveProfileImage,
		NeedToken:    true,
	},
	//cv_type  routes
	{
		HttpMethod:   http.MethodPost,
		RelativePath: "/profile/cv",
		Handler:      cvtype.UploadCv,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodPut,
		RelativePath: "/profile/cv",
		Handler:      cvtype.UpdateProfileCv,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodGet,
		RelativePath: "/profile/cv",
		Handler:      cvtype.GetProfileCv,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodGet,
		RelativePath: "/profile/cv/thumb",
		Handler:      cvtype.GetProfileCvThumb,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodDelete,
		RelativePath: "/profile/cv",
		Handler:      cvtype.RemoveProfileCv,
		NeedToken:    true,
	},
	//videos  routes
	{
		HttpMethod:   http.MethodPost,
		RelativePath: "/profile/video",
		Handler:      video.UploadVideo,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodPut,
		RelativePath: "/profile/video",
		Handler:      video.UpdateProfileVideo,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodGet,
		RelativePath: "/profile/video",
		Handler:      video.GetProfileVideo,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodDelete,
		RelativePath: "/profile/video",
		Handler:      video.RemoveProfileVideo,
		NeedToken:    true,
	},
	//Qr code authentication
	{
		HttpMethod:   http.MethodPost,
		RelativePath: "/generate-qr",
		Handler:      authentication.GenerateQrCode,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodPut,
		RelativePath: "/login/with-qr/:xid",
		Handler:      authentication.LoginWithQr,
		NeedToken:    false,
	},
	//Calendar planning routes
	{
		HttpMethod:   http.MethodPost,
		RelativePath: "/calendar",
		Handler:      planning.CreateUserPlannings,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodGet,
		RelativePath: "/calendar",
		Handler:      planning.GetUserPlannings,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodDelete,
		RelativePath: "/calendar",
		Handler:      planning.RemoveUserPlannings,
		NeedToken:    true,
	},
	//Calendar user planning routes
	{
		HttpMethod:   http.MethodPost,
		RelativePath: "/calendar/:calendar_id/:actor",
		Handler:      planning.AddUserIntoPlanning,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodGet,
		RelativePath: "/calendar/:calendar_id/actor",
		Handler:      planning.GetPlanningActors,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodDelete,
		RelativePath: "/calendar/:calendar_id/actor",
		Handler:      planning.RemoveUserFromPlanning,
		NeedToken:    true,
	},
	//Education Routes
	{
		HttpMethod:   http.MethodGet,
		RelativePath: "/education",
		Handler:      education.GetEducation,
		NeedToken:    false,
	},
	{
		HttpMethod:   http.MethodGet,
		RelativePath: "/education/:edu",
		Handler:      education.GetSubjects,
		NeedToken:    false,
	},
	//Education Level Routes'
	{
		HttpMethod:   http.MethodPost,
		RelativePath: "/user/education/",
		Handler:      education.SetUserEducationLevel,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodGet,
		RelativePath: "/user/education",
		Handler:      education.GetUserEducationLevel,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodPut,
		RelativePath: "/user/education/",
		Handler:      education.UpdateUserEducationLevel,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodGet,
		RelativePath: "/user/subject",
		Handler:      education.GetUserSubjects,
		NeedToken:    true,
	},

	//user_mark Routes
	{
		HttpMethod:   http.MethodPost,
		RelativePath: "/user_mark",
		Handler:      mark.RateUser,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodGet,
		RelativePath: "/user_mark/:userId",
		Handler:      mark.GetUserAverageMark,
		NeedToken:    true,
	},
	{
		HttpMethod:   http.MethodGet,
		RelativePath: "/user_mark/comment",
		Handler:      mark.GetUserMarkComment,
		NeedToken:    true,
	},
}
