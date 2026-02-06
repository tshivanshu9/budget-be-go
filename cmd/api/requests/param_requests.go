package requests

type IDParamRequest struct {
	Id uint `param:"id" validate:"required"`
}
