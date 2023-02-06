/*
 * @Author: SpenserCai
 * @Date: 2023-02-06 10:14:01
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-06 14:15:46
 * @Description: file content
 */
package lgscore

import (
	"errors"
	"os"
	"strings"

	"github.com/andygrunwald/vdf"
)

type SteamGame struct {
	GameName        string
	GameInstallPath string
}

type SteamApp struct {
	AppId         string
	PfxPath       string
	ProtonPath    string
	ProtonVersion string
	Game          SteamGame
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

func (s *SteamApp) SetGeProtonPath() error {
	// 判断s.PfxPath是否有值
	if s.PfxPath == "" {
		return errors.New("PfxPath is empty")
	}
	steamdataPath := strings.Split(s.PfxPath, "/pfx")[0]
	// 读取steamdataPath+"/config_info"文件
	configInfo, err := os.ReadFile(steamdataPath + "/config_info")
	if err != nil {
		return err
	}
	// 第一行是版本号，第二行是proton路径
	protonVersion := strings.Split(string(configInfo), "\n")[0]
	protonPath := strings.Split(string(configInfo), "\n")[1]
	// 去除第二行protonVersion后面的部分
	s.ProtonPath = strings.Split(protonPath, protonVersion)[0] + protonVersion + "/proton"
	s.ProtonVersion = protonVersion
	return nil

}

func (s *SteamApp) InitSteamApp() error {
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
				if _, ok := value.(map[string]interface{})["apps"].(map[string]interface{})[s.AppId]; ok {
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
					tmpPfxPath := value.(map[string]interface{})["path"].(string) + "/steamapps/compatdata/" + s.AppId + "/pfx"
					// 判断pfxPath目录是否存在
					if _, err := os.Stat(tmpPfxPath); err != nil {
						for _, localPath := range SteamPath {
							tmpPfxPath = HomePath + "/" + localPath + "/steamapps/compatdata/" + s.AppId + "/pfx"
							if _, err := os.Stat(tmpPfxPath); err == nil {
								s.PfxPath = tmpPfxPath
								break
							}
						}
					} else {
						s.PfxPath = tmpPfxPath
					}
					err = s.SetGeProtonPath()
					if err != nil {
						return err
					}
					return nil
				}
			}

		}
	}
	return nil
}
