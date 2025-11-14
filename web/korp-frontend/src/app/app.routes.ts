import { Routes } from '@angular/router';
import { ProductListComponent } from './components/products/product-list/product-list.component';
import { ProductFormComponent } from './components/products/product-form/product-form.component';
import { InvoiceListComponent } from './components/invoices/invoice-list/invoice-list.component';
import { InvoiceFormComponent } from './components/invoices/invoice-form/invoice-form.component';
import { InvoicePrintComponent } from './components/invoices/invoice-print/invoice-print.component';

export const routes: Routes = [
  // Rota padrão redireciona para produtos
  { path: '', redirectTo: '/products', pathMatch: 'full' },
  
  // Rotas de Produtos
  { path: 'products', component: ProductListComponent },
  { path: 'products/new', component: ProductFormComponent },
  
  // Rotas de Notas Fiscais
  { path: 'invoices', component: InvoiceListComponent },
  { path: 'invoices/new', component: InvoiceFormComponent },
  { path: 'invoices/:id/print', component: InvoicePrintComponent },
  
  // Rota 404
  { path: '**', redirectTo: '/products' }
];