INSERT INTO sources (name, api_name, base_url, enabled)
VALUES
  ('GDELT', 'GDELT API', 'https://api.gdeltproject.org/api/v2/doc/doc', TRUE),
  ('The Guardian', 'Guardian Open Platform', 'https://content.guardianapis.com', TRUE),
  ('Hacker News', 'Hacker News Firebase API', 'https://hacker-news.firebaseio.com/v0', TRUE),
  ('arXiv', 'arXiv API', 'http://export.arxiv.org/api/query', TRUE),
  ('Spaceflight News', 'Spaceflight News API', 'https://api.spaceflightnewsapi.net/v4', TRUE)
ON CONFLICT (name) DO NOTHING;

INSERT INTO evaluation_questions (question, expected_source, expected_keywords)
VALUES
  ('What are recent trends in AI infrastructure?', 'arXiv', ARRAY['AI infrastructure', 'large language models', 'inference', 'training']),
  ('Summarize recent space technology news.', 'Spaceflight News', ARRAY['space', 'launch', 'satellite', 'mission']),
  ('What are people discussing about Kubernetes or cloud infrastructure?', 'Hacker News', ARRAY['Kubernetes', 'cloud', 'infrastructure', 'discussion']),
  ('Find recent research related to large language models.', 'arXiv', ARRAY['large language models', 'transformers', 'research', 'abstract']),
  ('Compare recent news coverage and research discussion on cybersecurity.', 'GDELT', ARRAY['cybersecurity', 'news', 'threat', 'attack']),
  ('What is the latest coverage on renewable energy policy?', 'The Guardian', ARRAY['renewable energy', 'policy', 'climate', 'reporting']),
  ('What open-source tools are developers discussing for observability?', 'Hacker News', ARRAY['observability', 'open source', 'monitoring', 'developers']),
  ('Which papers mention retrieval augmented generation?', 'arXiv', ARRAY['retrieval augmented generation', 'RAG', 'embeddings', 'semantic search']),
  ('What are the most recent discussions around AI safety?', 'GDELT', ARRAY['AI safety', 'risk', 'policy', 'coverage']),
  ('What did readers discuss about startup funding and AI platforms?', 'Hacker News', ARRAY['startup', 'funding', 'AI platform', 'discussion']),
  ('Summarize Guardian articles about climate technology.', 'The Guardian', ARRAY['climate', 'technology', 'article', 'analysis']),
  ('Find space-related articles mentioning satellites or orbital launches.', 'Spaceflight News', ARRAY['satellite', 'launch', 'orbital', 'spacecraft']),
  ('What are recent global news items about semiconductors?', 'GDELT', ARRAY['semiconductor', 'supply chain', 'chip', 'news']),
  ('Find research discussing reinforcement learning advances.', 'arXiv', ARRAY['reinforcement learning', 'policy gradient', 'paper', 'research']),
  ('What tech discussions mention PostgreSQL or pgvector?', 'Hacker News', ARRAY['PostgreSQL', 'pgvector', 'database', 'semantic search']),
  ('What Guardian reporting covers health or biotech innovation?', 'The Guardian', ARRAY['health', 'biotech', 'innovation', 'reporting']),
  ('What news coverage mentions cloud security incidents?', 'GDELT', ARRAY['cloud security', 'incident', 'breach', 'coverage']),
  ('What are people saying about generative AI tooling?', 'Hacker News', ARRAY['generative AI', 'tooling', 'workflow', 'discussion']),
  ('Find recent studies about multimodal models.', 'arXiv', ARRAY['multimodal', 'vision-language', 'paper', 'study']),
  ('Summarize recent space science or exploration stories.', 'Spaceflight News', ARRAY['space science', 'exploration', 'mission', 'article']);
