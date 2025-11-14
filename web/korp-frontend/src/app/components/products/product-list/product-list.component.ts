import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { Subscription } from 'rxjs';
import { ProductService } from '../../../services/product.service';
import { Product } from '../../../models/product.model';

@Component({
  selector: 'app-product-list',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatTableModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
    MatProgressSpinnerModule,
    MatSnackBarModule
  ],
  templateUrl: './product-list.component.html',
  styleUrls: ['./product-list.component.scss']
})
export class ProductListComponent implements OnInit, OnDestroy {
  products: Product[] = [];
  loading = true;
  error: string | null = null;
  displayedColumns: string[] = ['code', 'description', 'balance', 'actions'];
  
  private subscription?: Subscription;

  constructor(
    private productService: ProductService,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.loadProducts();
  }

  ngOnDestroy(): void {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
  }

  loadProducts(): void {
    this.loading = true;
    this.error = null;
    
    this.subscription = this.productService.getProducts().subscribe({
      next: (products) => {
        this.products = products;
        this.loading = false;
      },
      error: (error) => {
        this.error = error.message;
        this.loading = false;
        console.error('Erro ao carregar produtos:', error);
      }
    });
  }

  deleteProduct(product: Product): void {
    if (!confirm(`Tem certeza que deseja excluir o produto ${product.code} - ${product.description}?`)) {
      return;
    }

    this.loading = true;
    this.productService.deleteProduct(product.id).subscribe({
      next: () => {
        this.showSuccess('Produto excluído com sucesso!');
        this.loadProducts();
      },
      error: (error) => {
        this.showError(error.message);
        this.loading = false;
      }
    });
  }

  refresh(): void {
    this.loadProducts();
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