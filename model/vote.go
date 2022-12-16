package model

type VoteParam struct {
	PostID    string `json:"post_id" binding:"required"`
	Direction int    `json:"direction,string" binding:"required,oneof=1 0 -1"` // 投票的方向，上下分别表示赞成反对，三种状态表示 1 0 -1
}

//func (v *VoteParam) UnmarshalJSON(data []byte) (err error) {
//	required := struct {
//		PostID    string `json:"post_id"`
//		Direction int    `json:"direction"`
//	}{}
//	err = json.Unmarshal(data, &required)
//	if err != nil {
//		return
//	} else if len(required.PostID) == 0 {
//		err = errors.New("缺少必填字段post_id")
//	} else if required.Direction == 0 {
//		err = errors.New("缺少必填字段direction")
//	} else {
//		v.PostID = required.PostID
//		v.Direction = required.Direction
//	}
//	return
//}
