package controller

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"testing"
)

func TestRequest_GetUserId(t *testing.T) {
	type fields struct {
		GCtx *gin.Context
		Val  string
	}
	ctxWithUid := &gin.Context{}
	ctxWithUid.Set("user_id", uint32(9))
	ctxWithoutUid := &gin.Context{}

	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			"had_set_uid",
			fields{
				GCtx: ctxWithUid,
			},
			9,
		},
		{
			"no_set_uid",
			fields{
				GCtx: ctxWithoutUid,
			},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Request{
				GCtx: tt.fields.GCtx,
				Val:  tt.fields.Val,
			}
			if got := req.GetUserId(); got != tt.want {
				t.Errorf("GetUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRequest(t *testing.T) {
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name string
		args args
		want Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRequest(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Bool(t *testing.T) {
	type fields struct {
		GCtx *gin.Context
		Val  string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Request{
				GCtx: tt.fields.GCtx,
				Val:  tt.fields.Val,
			}
			if got := req.Bool(); got != tt.want {
				t.Errorf("Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_GetAllParamsFromUrl(t *testing.T) {
	type fields struct {
		GCtx *gin.Context
		Val  string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Request{
				GCtx: tt.fields.GCtx,
				Val:  tt.fields.Val,
			}
			if got := req.GetAllParamsFromUrl(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllParamsFromUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_GetHeader(t *testing.T) {
	type fields struct {
		GCtx *gin.Context
		Val  string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Request{
				GCtx: tt.fields.GCtx,
				Val:  tt.fields.Val,
			}
			if got := req.GetHeader(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_GetQuery(t *testing.T) {
	type fields struct {
		GCtx *gin.Context
		Val  string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Request{
				GCtx: tt.fields.GCtx,
				Val:  tt.fields.Val,
			}
			if got := req.GetQuery(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_GetQueryDefault(t *testing.T) {
	type fields struct {
		GCtx *gin.Context
		Val  string
	}
	type args struct {
		key          string
		defaultValue string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Request{
				GCtx: tt.fields.GCtx,
				Val:  tt.fields.Val,
			}
			if got := req.GetQueryDefault(tt.args.key, tt.args.defaultValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetQueryDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_IsMimeJson(t *testing.T) {
	type fields struct {
		GCtx *gin.Context
		Val  string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Request{
				GCtx: tt.fields.GCtx,
				Val:  tt.fields.Val,
			}
			if got := req.IsMimeJson(); got != tt.want {
				t.Errorf("IsMimeJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_PostForm(t *testing.T) {
	type fields struct {
		GCtx *gin.Context
		Val  string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Request{
				GCtx: tt.fields.GCtx,
				Val:  tt.fields.Val,
			}
			if got := req.PostForm(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostForm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_PostToModel(t *testing.T) {
	type fields struct {
		GCtx *gin.Context
		Val  string
	}
	type args struct {
		obj any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Request{
				GCtx: tt.fields.GCtx,
				Val:  tt.fields.Val,
			}
			if got := req.PostToModel(tt.args.obj); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostToModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_ShouldBindJSON(t *testing.T) {
	type fields struct {
		GCtx *gin.Context
		Val  string
	}
	type args struct {
		obj any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Request{
				GCtx: tt.fields.GCtx,
				Val:  tt.fields.Val,
			}
			req.ShouldBindJSON(tt.args.obj)
		})
	}
}

func TestRequest_ShouldBindQuery(t *testing.T) {
	type fields struct {
		GCtx *gin.Context
		Val  string
	}
	type args struct {
		obj any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Request{
				GCtx: tt.fields.GCtx,
				Val:  tt.fields.Val,
			}
			req.ShouldBindQuery(tt.args.obj)
		})
	}
}

func TestRequest_Value(t *testing.T) {
	type fields struct {
		GCtx *gin.Context
		Val  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Request{
				GCtx: tt.fields.GCtx,
				Val:  tt.fields.Val,
			}
			if got := req.Value(); got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
