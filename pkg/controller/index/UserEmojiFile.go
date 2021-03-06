package index

import (
	"emoji/pkg/config"
	"emoji/pkg/database"
	"emoji/pkg/model/entity"
	"emoji/pkg/model/logic"
	"emoji/pkg/system"
	"emoji/pkg/unity"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"html"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type UserEmojiFile struct {
	userUniqueId    string
	userAssFilePath string
	userFileSave    string
	sysFilePath     string
	sysAssPath      string
}

func NewUserEmojiFile()*UserEmojiFile  {
	return &UserEmojiFile{}
}

func (this *UserEmojiFile) EmojiGenerator(ctx *gin.Context){
	var emoji     entity.EmojiFile
	var userEmoji entity.UserEmojiFile
	var err       error
	emoji.Md5Encode  = html.EscapeString(ctx.Query("encode"))
	body,_ := ioutil.ReadAll(ctx.Request.Body)
	err     = json.Unmarshal(body,&emoji)
	protocol,_ := ctx.Get("protocol")
	unity.ErrorCheck(err)
	sysFileInfo := logic.NewSysEmojiFileLogic(database.GetOrm()).GetSysFileListFirst(emoji)
	if len(sysFileInfo) != 0 {
		this.sysAssPath = sysFileInfo["base_path"].(string) + sysFileInfo["name"].(string) + ".ass"
		this.sysFilePath= sysFileInfo["path"].(string)
		if userId,ok := ctx.Get("openId");ok{
			userEmoji.CreateTime = unity.GetNowDateTime(config.SecondFormat)
			this.userAssFilePath,err = this.AnalysisAss(sysFileInfo["name"].(string),emoji.Sentence,userId.(string))
			if err == nil {
				this.userUniqueId = userId.(string)
				userEmoji.OpenId  = userId.(string)
				userEmoji.ImageUrl= protocol.(string) + ctx.Request.Host + this.userFileSave + ".gif"
				if this.ExecuteCommand(){
					if logic.NewUserEmojiFileLogic(database.GetOrm()).InsertNewRecord(userEmoji){
						system.PrintSuccess(ctx,201,"",map[string]interface{}{
							"image_url" : protocol.(string) + ctx.Request.Host + this.userFileSave + ".gif",
						})
						return
					}else{
						system.PrintException(ctx,223,"",map[string]interface{}{})
						return
					}
				}else{
					system.PrintException(ctx,222,"",map[string]interface{}{})
					return
				}
			}else{
				system.PrintException(ctx,220,"",map[string]interface{}{})
				return
			}
		}
		system.PrintException(ctx,401,"",map[string]interface{}{})
		return
	}
	system.PrintException(ctx,112,"",map[string]interface{}{})
}

func (this *UserEmojiFile)AnalysisAss(fileName string,sentence string,userId string)(string,error)  {
	var fileNewStr string
	pkgPrefix   := "./pkg"
	fileStr,err := ioutil.ReadFile(this.sysAssPath)
	fileNewStr   = string(fileStr)
	regRule := "Dialogue: (\\d,\\d:\\d{0,2}:\\d{0,2}\\.\\d{0,2}){2},\\w+,(,\\d{0,2}){3}(,){2}<\\?loading-%s\\?>";
	unity.ErrorCheck(err)
	sentenceStr:= strings.Split(sentence,"|")
	for key,val := range sentenceStr{
		matchLineReg,err := regexp.Compile(fmt.Sprintf(regRule,strconv.Itoa(key)))
		unity.ErrorCheck(err)
		matchLineStr := matchLineReg.FindString(string(fileStr))
		if len(matchLineStr) == 0 {
			return "",errors.New("error")
		}
		matchLinePartReg,err := regexp.Compile("Dialogue: (\\d,\\d:\\d{0,2}:\\d{0,2}\\.\\d{0,2}){2},\\w+,(,\\d{0,2}){3}(,){2}")
		unity.ErrorCheck(err)
		matchLinePartStr := matchLinePartReg.FindString(matchLineStr)
		matchLinePartStr += val
		fileNewStr = matchLineReg.ReplaceAllString(fileNewStr,matchLinePartStr)
	}
	userNewFile    := "/assets/user" + config.OS_SEPREATOR
	userNewFile    += unity.GetNowDateTime(config.HourFormat) + config.OS_SEPREATOR
	this.userFileSave = userNewFile + unity.Md5String(this.userUniqueId + unity.GetNowDateTime(config.HourFormat))
	if unity.DirExistValidate(pkgPrefix + userNewFile) == false{
		unity.DirMakeAll(pkgPrefix + userNewFile)
	}
	userNewFile    = pkgPrefix + userNewFile
	userNewFile    += unity.Md5String(this.userUniqueId + unity.GetNowDateTime(config.HourFormat))
	userNewFile    += config.ASS_FILE_EXT
	unity.FileMake(userNewFile)
	err = ioutil.WriteFile(userNewFile,[]byte(fileNewStr),os.ModePerm)
	unity.ErrorCheck(err)
	return userNewFile,nil
}

func (this *UserEmojiFile)ExecuteCommand() (bool) {
	var command =  &exec.Cmd{}
	pkgPrefix   := "./pkg"
	gifSuffix   := ".gif"
	sysFilePath := this.sysFilePath
	usrAssFile  := this.userAssFilePath
	usrSavePath := pkgPrefix  + this.userFileSave + gifSuffix
	command      = exec.Command("ffmpeg","-y","-i",sysFilePath,"-vf",fmt.Sprintf("ass=%s",usrAssFile),"-r","10","-b:v","1500k","-s","250*200","-bufsize","1500k",usrSavePath)
	if _,err := command.CombinedOutput();err != nil{
		unity.ErrorCheck(err)
	}
	if err := os.Remove(usrAssFile);err != nil {

	}
	return true
}
