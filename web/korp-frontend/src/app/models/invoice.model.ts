// Model de Nota Fiscal
export interface Invoice {
  id: string;
  number: number;
  status: InvoiceStatus;
  items: InvoiceItem[];
  created_at: string;
  updated_at: string;
  closed_at?: string;
}

// Status da Nota Fiscal
export enum InvoiceStatus {
  OPEN = 'ABERTA',
  CLOSED = 'FECHADA'
}

// Item da Nota Fiscal
export interface InvoiceItem {
  product_id: string;
  product_code: string;
  description: string;
  quantity: number;
}

// DTO para criação de nota fiscal
export interface CreateInvoiceDTO {
  items: CreateInvoiceItemDTO[];
}

// DTO para item ao criar nota
export interface CreateInvoiceItemDTO {
  product_id: string;
  quantity: number;
}

// Resposta de impressão
export interface PrintResponse {
  success: boolean;
  message: string;
  invoice?: Invoice;
}