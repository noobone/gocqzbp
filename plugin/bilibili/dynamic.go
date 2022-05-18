package bilibili

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	control "github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/process"
)

const (
	BILIBILI    = "bilibili"
	BILIBILI_CN = "哔哩哔哩"
	BILIBILI_DB = BILIBILI + ".db"
)

func init() {
	engine := control.Register(BILIBILI, &control.Options{
		DisableOnDefault: true,
		Help:             BILIBILI_CN,
		PublicDataFolder: BILIBILI,
	})

	// getdb := ctxext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
	// 	db.DBPath = engine.DataFolder() + BILIBILI_DB
	// 	_, err := engine.GetLazyData(BILIBILI_DB, true)
	// 	if err != nil {
	// 		ctx.SendChain(message.Text("ERROR:", err))
	// 		return false
	// 	}
	// 	err = db.Create(BILIBILI, &curse{})
	// 	if err != nil {
	// 		ctx.SendChain(message.Text("ERROR:", err))
	// 		return false
	// 	}
	// 	cnt, err := db.Count(BILIBILI)
	// 	if err != nil {
	// 		ctx.SendChain(message.Text("ERROR:", err))
	// 		return false
	// 	}
	// 	logrus.Infof("[%s]加载 %d 条数据", BILIBILI, cnt)
	// 	return true
	// })

	engine.OnFullMatch("骂我", getdb).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		process.SleepAbout1sTo2s()
		text := getRandomCurseByLevel(minLevel).Text
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(text))
	})

	engine.OnFullMatch("大力骂我", getdb).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		process.SleepAbout1sTo2s()
		text := getRandomCurseByLevel(maxLevel).Text
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(text))
	})

	engine.OnKeywordGroup([]string{"他妈", "公交车", "你妈", "操", "屎", "去死", "快死", "我日", "逼", "尼玛", "艾滋", "癌症", "有病", "烦你", "你爹", "屮", "cnm"}, zero.OnlyToMe, getdb).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			text := getRandomCurseByLevel(maxLevel).Text
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(text))
		})
}
