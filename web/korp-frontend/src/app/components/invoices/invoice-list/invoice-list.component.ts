import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatChipsModule } from '@angular/material/chips';
import { Subscription } from 'rxjs';
import { InvoiceService } from '../../../services/invoice.service';
import { Invoice, InvoiceStatus } from '../../../models/invoice.model';

@Component({
  selector: 'app-invoice-list',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatTableModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
    MatProgressSpinnerModule,
    MatChipsModule
  ],
  templateUrl: './invoice-list.component.html',
  styleUrls: ['./invoice-list.component.scss']
})
export class InvoiceListComponent implements OnInit, OnDestroy {
  invoices: Invoice[] = [];
  loading = true;
  error: string | null = null;
  displayedColumns: string[] = ['number', 'status', 'items', 'created_at', 'actions'];
  InvoiceStatus = InvoiceStatus; // Para usar no template
  
  private subscription?: Subscription;

  constructor(private invoiceService: InvoiceService) {}

  ngOnInit(): void {
    this.loadInvoices();
  }

  ngOnDestroy(): void {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
  }

  /**
   * Carrega a lista de notas fiscais
   */
  loadInvoices(): void {
    this.loading = true;
    this.error = null;
    
    this.subscription = this.invoiceService.getInvoices().subscribe({
      next: (invoices) => {
        this.invoices = invoices;
        this.loading = false;
      },
      error: (error) => {
        this.error = error.message;
        this.loading = false;
        console.error('Erro ao carregar notas fiscais:', error);
      }
    });
  }

  /**
   * Recarrega a lista de notas fiscais
   */
  refresh(): void {
    this.loadInvoices();
  }

  /**
   * Retorna a quantidade total de itens de uma nota
   */
  getTotalItems(invoice: Invoice): number {
    return invoice.items.reduce((sum, item) => sum + item.quantity, 0);
  }

  /**
   * Formata data para exibição
   */
  formatDate(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleString('pt-BR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
}