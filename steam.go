/*
 * @Author: SpenserCai
 * @Date: 2023-02-06 10:14:01
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-06 11:07:40
 * @Description: file content
 */
package lgscore

import (
	"os"

	"github.com/andygrunwald/vdf"
)

type SteamGame struct {
	GameName        string
	GameInstallPath string
}

type SteamApp struct {
	AppId      string
	PfxPath    string
	ProtonPath string
	Game       SteamGame
}

func (s *SteamApp) SetGameInfo(appsPath string) error {
	f, err := os.Open(appsPath + "/appmanifest_" + s.AppId + ".acf")
	if err != nil {
		return err
	}
	defer f.Close()
	p := vdf.NewParser(f)
	v, err := p.Parse()
	if err != nil {
		return err
	}
	s.Game.GameName = v["AppState"].(map[string]interface{})["name"].(string)
	s.Game.GameInstallPath = v["AppState"].(map[string]interface{})["installdir"].(string)
	return nil
}

func (s *SteamApp) InitSteamApp(appid string) error {
	for i := 0; i < len(SteamPath); i++ {
		appPath := HomePath + "/" + SteamPath[i] + "/steamapps"
		// 判断路径是否存在libaryfolders.vdf文件
		if _, err := os.Stat(appPath + "/libraryfolders.vdf"); err == nil {
			// 读取libaryfolders.vdf文件，遍历第一层
			f, err := os.Open(appPath + "/libraryfolders.vdf")
			if err != nil {
				return err
			}
			defer f.Close()
			p := vdf.NewParser(f)
			v, err := p.Parse()
			if err != nil {
				return err
			}
			// fmt.Println(v)
			// 遍历map->libraryfolders
			for _, value := range v["libraryfolders"].(map[string]interface{}) {
				// 判断map的apps中是否有名为appid的键，如果存在则设置steamApp的GamePath为map的path+"/steamapps/common/"+gameName
				if _, ok := value.(map[string]interface{})["apps"].(map[string]interface{})[appid]; ok {
					appsPath := value.(map[string]interface{})["path"].(string) + "/steamapps"
					err := s.SetGameInfo(appsPath)
					if err != nil {
						return err
					}
					tmpGamePath := value.(map[string]interface{})["path"].(string) + "/steamapps/common/" + s.Game.GameInstallPath
					// 判断GamePath目录是否存在
					if _, err := os.Stat(tmpGamePath); err != nil {
						return err
					}
					s.Game.GameInstallPath = tmpGamePath
					s.PfxPath = value.(map[string]interface{})["path"].(string) + "/steamapps/compatdata/" + appid + "/pfx"
					return nil
				}
			}

		}
	}
	return nil
}
