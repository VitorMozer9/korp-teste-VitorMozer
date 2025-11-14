import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
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
    MatProgressSpinnerModule
  ],
  templateUrl: './product-list.component.html',
  styleUrls: ['./product-list.component.scss']
})
export class ProductListComponent implements OnInit, OnDestroy {
  products: Product[] = [];
  loading = true;
  error: string | null = null;
  displayedColumns: string[] = ['code', 'description', 'balance'];
  
  private subscription?: Subscription;

  constructor(private productService: ProductService) {}

  /**
   * ngOnInit - Ciclo de vida do Angular
   * Executado quando o componente é inicializado
   */
  ngOnInit(): void {
    this.loadProducts();
  }

  /**
   * ngOnDestroy - Ciclo de vida do Angular
   * Executado quando o componente é destruído
   * IMPORTANTE: Limpar subscriptions para evitar memory leaks
   */
  ngOnDestroy(): void {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
  }

  /**
   * Carrega a lista de produtos
   * Usa RxJS Observable do serviço
   */
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

  /**
   * Recarrega a lista de produtos
   */
  refresh(): void {
    this.loadProducts();
  }
}