
resource "hpe_opsramp_kb_article" "kb_article_default" {
  subject     = "Default article"
  content     = "Default article's description"
  category_id = hpe_opsramp_kb_category.kb_category_default.id
}