// Model de Produto
export interface Product {
  id: string;
  code: string;
  description: string;
  balance: number;
  created_at: string;
  updated_at: string;
}

// DTO para criação de produto
export interface CreateProductDTO {
  code: string;
  description: string;
  balance: number;
}

// DTO para atualização de produto
export interface UpdateProductDTO {
  code: string;
  description: string;
  balance: number;
}