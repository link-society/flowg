use serde_json::{json, Value};

use super::lexer::{Token, TokenStream};

peg::parser!{
  pub grammar parser() for TokenStream {
    pub rule start() -> String
      = expr:expr()
      { expr.to_string() }

    rule expr() -> Value
      = precedence!{
        lhs:(@) [Token::OpOr] rhs:@ {
          json!({"$or": [lhs, rhs]})
        }
        --
        lhs:(@) [Token::OpAnd] rhs:@ {
          json!({"$and": [lhs, rhs]})
        }
        --
        [Token::OpNot] expr:@ {
          json!({"$not": expr})
        }
        --
        t:field_eq() { t }
        t:field_ne() { t }
        t:field_in() { t }
        t:field_nin() { t }
        --
        [Token::LParen] expr:expr() [Token::RParen] { expr }
      }

    rule field_eq() -> Value
      = field:field() [Token::Eq] value:string()
      { json!({"$eq": {"field": field, "value": value}}) }

    rule field_ne() -> Value
      = field:field() [Token::Ne] value:string()
      { json!({"$not": {"$eq": {"field": field, "value": value}}}) }

    rule field_in() -> Value
      = field:field() [Token::OpIn] [Token::LBracket] values:(string() ** [Token::Comma]) [Token::RBracket]
      { json!({"$in": {"field": field, "values": values}}) }

    rule field_nin() -> Value
      = field:field() [Token::OpNot] [Token::OpIn] [Token::LBracket] values:(string() ** [Token::Comma]) [Token::RBracket]
      { json!({"$not": {"$in": {"field": field, "values": values}}}) }

    rule field() -> Value
      = [Token::Ident(s)]
      { json!{s} }

    rule string() -> Value
      = [Token::String(s)]
      { json!{s} }
  }
}
