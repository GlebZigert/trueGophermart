package accrual

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"time"

	"github.com/GlebZigert/trueGophermart/internal/config"
	"github.com/GlebZigert/trueGophermart/internal/logger"
	"github.com/GlebZigert/trueGophermart/internal/model"
	"gorm.io/gorm"
)

//подобрать флаг адреса акруала

//если флага нет - выставить ошибку

// соедниниться с акруалом - соединения нет - ошибка?

//соединение пропало - ошибка?

// периодически с интервалом ( задан флагом GET_ORDER_FROM_DB_TO_GET_ACCRUAL) выполнять следующие действия:

//забирать значения из бд - нбирать определенный лимит (задан флагом)

// посылать эти значения в аккруал

//при получении ответа производить запись инфы по этому ордру в бд и убирать ордер из внутр

type Accrual struct {
	DB     *gorm.DB
	cfg    *config.Config
	logger logger.Logger
}

func NewAccrual(db *gorm.DB, cfg *config.Config, logger logger.Logger, ctx context.Context) (*Accrual, error) {

	aq := &Accrual{db, cfg, logger}

	go aq.Run(ctx)

	return aq, nil
}

func (aq *Accrual) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	aq.logger.Info("accrual --> ", nil)
	for {

		select {

		case <-ticker.C:
			//берем ордеры из БД

			var orders []model.Order
			result := aq.DB.Where("Status=?", model.ORDER_REGISTERED).Find(&orders)

			if result.Error != nil {
				aq.logger.Error("accrual : ", map[string]interface{}{
					"err": result.Error,
				})
				continue

			}

			if len(orders) == 0 {
				continue
			}

			aq.logger.Info("Взял в обработку : ", map[string]interface{}{
				"orders": orders,
			})

			for _, order := range orders {
				req := aq.cfg.AccrualAddress + "/api/orders/" + order.Number
				aq.logger.Info("accrual : ", map[string]interface{}{
					"req": req,
				})
				resp, err := http.Get(req)
				if err != nil {
					aq.logger.Error("accrual http.Get(req): ", map[string]interface{}{
						"err": err.Error(),
					})
					continue
				}

				aq.logger.Info("accrual : ", map[string]interface{}{
					"resp": resp,
				})

				if resp.StatusCode != http.StatusOK {
					continue
				}

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					aq.logger.Error("accrual io.ReadAll: ", map[string]interface{}{
						"err": err.Error(),
					})
					continue

				}

				var answer model.Answer

				err = json.Unmarshal(body, &answer)

				if err != nil {
					aq.logger.Error("accrual json.Unmarshal: ", map[string]interface{}{
						"err": err.Error(),
					})
					continue

				}

				aq.logger.Info("answer : ", map[string]interface{}{
					"answer":  answer,
					"Accrual": answer.Accrual,
					"Status":  answer.Status,
					"user":    order.UID,
				})

				order.Accrual = answer.Accrual
				order.Status = answer.Status

				aq.DB.Save(order)

				var user model.User
				//находим пользователя
				res := aq.DB.Where("id=?", order.UID).First(&user)

				if res.Error != nil {
					aq.logger.Error("err find user : ", map[string]interface{}{
						"err": res.Error,
					})
					continue
				}
				user.Current = user.Current + order.Accrual

				aq.logger.Info("user : ", map[string]interface{}{
					"user.ID":        user.ID,
					"user.Current":   user.Current,
					"user.withdrawn": user.withdrawn,
				})
				aq.DB.Save(user)
			}

		}
	}

}
