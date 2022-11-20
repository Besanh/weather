package service

import (
	dataUtil "weather/common/data"
)

type PBXInformation struct {
	PBXDomain, APIDomain, PBXPort, PBXWss, PBXOutboundProxy string
	PBXTransport                                            string
}

var LeadFileDir = "import/lead/"

var PBXInfo PBXInformation

type ByTime []map[string]interface{}

func (a ByTime) Len() int { return len(a) }
func (a ByTime) Less(i, j int) bool {
	iTimeStr := a[i]["time"].(string)
	jTimeStr := a[j]["time"].(string)
	iTime := dataUtil.ParseFromStringToTime(iTimeStr)
	jTime := dataUtil.ParseFromStringToTime(jTimeStr)
	return iTime.After(jTime)
}
func (a ByTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// func GetUserInfo(ctx context.Context, userUuid string) (*model.UserData, error) {
// 	user, err := repository.UserRepo.SelectUserDataById(ctx, userUuid)
// 	if err != nil {
// 		log.Error(err)
// 		return nil, err
// 	} else if user != nil {
// 		if dataUtil.InSlice(user.Level, []string{SUPERADMIN, ADMIN, MANAGER, LEADER}) {
// 			groupUuids := make([]string, 0)
// 			userUuids := make([]string, 0)
// 			if dataUtil.InSlice(user.Level, []string{MANAGER, LEADER}) {
// 				for _, group := range user.Groups {
// 					groupUuids = append(groupUuids, group.GroupUuid)
// 				}
// 			}
// 			if users, err := repository.UserRepo.SelectUsersInfoOfGroupUsers(ctx, user.DomainUuid, userUuid, user.Level, groupUuids, []string{}); err != nil {
// 				log.Error(err)
// 				return nil, errors.New("user info is invalid")
// 			} else if users != nil {
// 				user.ManageUsers = *users
// 				for _, u := range *users {
// 					userUuids = append(userUuids, u.UserUuid)
// 				}
// 			}
// 			if len(userUuids) > 0 {
// 				if extensions, err := repository.ExtensionRepo.SelectExtensionsOfUserUuids(ctx, user.DomainUuid, userUuids); err != nil {
// 					log.Error(err)
// 					return nil, errors.New("user info is invalid")
// 				} else if extensions != nil {
// 					user.ManageExtensions = *extensions
// 				}
// 			}
// 			return user, nil
// 		} else {
// 			return user, nil
// 		}
// 	} else {
// 		return nil, errors.New("user info is invalid")
// 	}
// }
