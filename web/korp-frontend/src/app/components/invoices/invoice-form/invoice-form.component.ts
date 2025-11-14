import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, FormArray, Validators, ReactiveFormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTableModule } from '@angular/material/table';
import { ProductService } from '../../../services/product.service';
import { InvoiceService } from '../../../services/invoice.service';
import { Product } from '../../../models/product.model';
import { CreateInvoiceDTO, CreateInvoiceItemDTO } from '../../../models/invoice.model';

@Component({
  selector: 'app-invoice-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatSelectModule,
    MatInputModule,
    MatButtonModule,
    MatCardModule,
    MatIconModule,
    MatSnackBarModule,
    MatTableModule
  ],
  templateUrl: './invoice-form.component.html',
  styleUrls: ['./invoice-form.component.scss']
})
export class InvoiceFormComponent implements OnInit {
  invoiceForm: FormGroup;
  products: Product[] = [];
  loading = false;
  loadingProducts = true;
  displayedColumns: string[] = ['product', 'quantity', 'actions'];

  constructor(
    private fb: FormBuilder,
    private productService: ProductService,
    private invoiceService: InvoiceService,
    private router: Router,
    private snackBar: MatSnackBar
  ) {
    this.invoiceForm = this.fb.group({
      items: this.fb.array([], Validators.required)
    });
  }

  ngOnInit(): void {
    this.loadProducts();
  }

  /**
   * Carrega lista de produtos disponíveis
   */
  loadProducts(): void {
    this.productService.getProducts().subscribe({
      next: (products) => {
        this.products = products.filter(p => p.balance > 0);
        this.loadingProducts = false;
        if (this.products.length > 0) {
          this.addItem(); // Adiciona primeira linha automaticamente
        }
      },
      error: (error) => {
        this.showError(error.message);
        this.loadingProducts = false;
      }
    });
  }

  /**
   * Retorna o FormArray de itens
   */
  get items(): FormArray {
    return this.invoiceForm.get('items') as FormArray;
  }

  /**
   * Cria um novo FormGroup para item
   */
  createItem(): FormGroup {
    return this.fb.group({
      product_id: ['', Validators.required],
      quantity: [1, [Validators.required, Validators.min(1)]]
    });
  }

  /**
   * Adiciona novo item à nota
   */
  addItem(): void {
    this.items.push(this.createItem());
  }

  /**
   * Remove item da nota
   */
  removeItem(index: number): void {
    this.items.removeAt(index);
  }

  /**
   * Retorna o produto selecionado em um item
   */
  getSelectedProduct(index: number): Product | undefined {
    const productId = this.items.at(index).get('product_id')?.value;
    return this.products.find(p => p.id === productId);
  }

  /**
   * Valida se a quantidade solicitada está disponível
   */
  validateQuantity(index: number): boolean {
    const item = this.items.at(index);
    const product = this.getSelectedProduct(index);
    const quantity = item.get('quantity')?.value;
    
    if (!product || !quantity) return true;
    return quantity <= product.balance;
  }

  /**
   * Submete o formulário
   */
  onSubmit(): void {
    if (this.invoiceForm.invalid) {
      this.showError('Por favor, preencha todos os campos corretamente');
      return;
    }

    // Valida quantidades
    for (let i = 0; i < this.items.length; i++) {
      if (!this.validateQuantity(i)) {
        this.showError(`Quantidade solicitada para o item ${i + 1} excede o saldo disponível`);
        return;
      }
    }

    this.loading = true;
    const invoice: CreateInvoiceDTO = {
      items: this.items.value as CreateInvoiceItemDTO[]
    };

    this.invoiceService.createInvoice(invoice).subscribe({
      next: () => {
        this.showSuccess('Nota fiscal criada com sucesso!');
        this.router.navigate(['/invoices']);
      },
      error: (error) => {
        this.showError(error.message);
        this.loading = false;
      }
    });
  }

  onCancel(): void {
    this.router.navigate(['/invoices']);
  }

  private showSuccess(message: string): void {
    this.snackBar.open(message, 'Fechar', {
      duration: 3000,
      horizontalPosition: 'end',
      verticalPosition: 'top',
      panelClass: ['success-snackbar']
    });
  }

  private showError(message: string): void {
    this.snackBar.open(message, 'Fechar', {
      duration: 5000,
      horizontalPosition: 'end',
      verticalPosition: 'top',
      panelClass: ['error-snackbar']
    });
  }
}