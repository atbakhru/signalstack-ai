package config

type OpenAIConfig struct {
    APIKey        string
    EmbeddingModel string
    ChatModel      string
}

func NewOpenAIConfig(cfg Config) OpenAIConfig {
    return OpenAIConfig{
        APIKey:        cfg.OpenAIAPIKey,
        EmbeddingModel: cfg.OpenAIEmbeddingModel,
        ChatModel:      cfg.OpenAIChatModel,
    }
}
