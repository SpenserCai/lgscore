/*
 * @Author: SpenserCai
 * @Date: 2023-02-06 09:47:07
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-06 10:19:44
 * @Description: file content
 */
package lgscore

import "os"

var HomePath = os.Getenv("HOME")

var SteamPath = []string{
	".steam/steam",
	".local/share/Steam",
}
