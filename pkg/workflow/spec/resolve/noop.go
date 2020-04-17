package resolve

var (
	NoOpDataTypeResolver      DataTypeResolver      = ChainDataTypeResolvers()
	NoOpSecretTypeResolver    SecretTypeResolver    = ChainSecretTypeResolvers()
	NoOpOutputTypeResolver    OutputTypeResolver    = ChainOutputTypeResolvers()
	NoOpParameterTypeResolver ParameterTypeResolver = ChainParameterTypeResolvers()
	NoOpAnswerTypeResolver    AnswerTypeResolver    = ChainAnswerTypeResolvers()
)
