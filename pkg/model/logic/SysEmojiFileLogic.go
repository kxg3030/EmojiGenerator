package logic

import (
	"emoji/pkg/model/entity"
	"emoji/pkg/unity"
	"github.com/gohouse/gorose"
	"strconv"
	"strings"
)

type SysEmojiFileLogic struct {
	Orm  *gorose.Session
}

func NewSysEmojiFileLogic(orm *gorose.Session)*SysEmojiFileLogic  {
	return &SysEmojiFileLogic{
		Orm : orm,
	}
}

// insert new record
func (this *SysEmojiFileLogic) InsertNewFileRecord(emoji entity.EmojiFile) bool  {
	result ,err  := this.Orm.Table("sys_emoji_file").Where(map[string]interface{}{
		"md5_encode"  : emoji.Md5Encode,
		"extension"   : emoji.Extension,
	}).First()
	if len(result) >= 1{
		return false
	}
	insertId,err := this.Orm.Table("sys_emoji_file").Data(map[string]interface{}{
		"name"        : emoji.Name,
		"path"        : emoji.Path,
		"extension"   : emoji.Extension,
		"md5_encode"  : emoji.Md5Encode,
		"create_time" : emoji.CreateTime,
		"base_path"   : emoji.BasePath,
		"sentence"    : emoji.Sentence,
		"image_url"   : emoji.ImageUrl,
		"cover_url"   : emoji.CoverUrl,
		"sentence_count"    : emoji.SentenceCount,
	}).InsertGetId()
	unity.ErrorCheck(err)
	return insertId >= 1
}

// get sys file from database
func (this *SysEmojiFileLogic)GetSysFileList(emoji entity.EmojiFile)[]map[string]interface{}  {
	result,err := this.Orm.Table("sys_emoji_file").Where(map[string]interface{}{
		"md5_encode" : emoji.Md5Encode,
	}).Fields("path,base_path,extension,name,sentence_count").Get()
	unity.ErrorCheck(err)
	return result
}

// get sys file from database
func (this *SysEmojiFileLogic)GetSysFileListFirst(emoji entity.EmojiFile)map[string]interface{}  {
	result,err := this.Orm.Table("sys_emoji_file").Where(map[string]interface{}{
		"md5_encode" : emoji.Md5Encode,
		"extension"  : ".mp4",
	}).Fields("path,base_path,extension,name").First()
	unity.ErrorCheck(err)
	return result
}

// update database column
func (this *SysEmojiFileLogic)UpdateSysFileImageUrl(url string,cover string,md5 string,count int64)bool  {
	_,err := this.Orm.Table("sys_emoji_file").Where("md5_encode",md5).Data(map[string]interface{}{
		"image_url" : url,
		"cover_url" : cover,
		"sentence_count" : count,
	}).Update()
	unity.ErrorCheck(err)
	return true
}

// select record
func (this *SysEmojiFileLogic)SelectSysFileList(filed ,page ,size string)[]map[string]interface{}  {
	pageNum,_ := strconv.Atoi(page)
	pageSize,_:= strconv.Atoi(size)
	startIndex := (pageNum - 1) * pageSize
	result,err := this.Orm.Table("sys_emoji_file").
		Where("extension",".ass").
		Fields(filed).
		Offset(startIndex).
		Limit(pageSize).
		Get()
	unity.ErrorCheck(err)
	if len(result) == 0{
		return []map[string]interface{}{}
	}
	return result
}

// get one record
func (this *SysEmojiFileLogic)GetsSysFileFirstById(id string)map[string]interface{}  {
	filed := "id,image_url,sentence_count,sentence,md5_encode"
	result,err := this.Orm.Table("sys_emoji_file").Where("id",id).Fields(filed).First()
	if err != nil {
		return map[string]interface{}{}
	}
	if len(result) > 1{
		for key ,value := range result{
			if key == "sentence"{
				result[key] = strings.Split(value.(string),"|")
			}
		}
	}
	return result
}