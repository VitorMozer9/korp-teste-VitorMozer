import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTableModule } from '@angular/material/table';
import { MatChipsModule } from '@angular/material/chips';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { Subscription } from 'rxjs';
import { InvoiceService } from '../../../services/invoice.service';
import { Invoice, InvoiceStatus } from '../../../models/invoice.model';

@Component({
  selector: 'app-invoice-print',
  standalone: true,
  imports: [
    CommonModule,
    MatButtonModule,
    MatCardModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatTableModule,
    MatChipsModule,
    MatDialogModule
  ],
  templateUrl: './invoice-print.component.html',
  styleUrls: ['./invoice-print.component.scss']
})
export class InvoicePrintComponent implements OnInit, OnDestroy {
  invoice: Invoice | null = null;
  loading = true;
  printing = false;
  error: string | null = null;
  InvoiceStatus = InvoiceStatus;
  displayedColumns: string[] = ['code', 'description', 'quantity'];
  
  private subscription?: Subscription;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private invoiceService: InvoiceService,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.loadInvoice(id);
    }
  }

  ngOnDestroy(): void {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
  }

  /**
   * Carrega os dados da nota fiscal
   */
  loadInvoice(id: string): void {
    this.loading = true;
    this.error = null;
    
    this.subscription = this.invoiceService.getInvoice(id).subscribe({
      next: (invoice) => {
        this.invoice = invoice;
        this.loading = false;
      },
      error: (error) => {
        this.error = error.message;
        this.loading = false;
      }
    });
  }

  /**
   * Imprime (fecha) a nota fiscal
   * Esta é a operação que atualiza o estoque
   */
  printInvoice(): void {
    if (!this.invoice || this.invoice.status !== InvoiceStatus.OPEN) {
      this.showError('Esta nota não pode ser impressa');
      return;
    }

    this.printing = true;
    
    this.invoiceService.printInvoice(this.invoice.id).subscribe({
      next: (response) => {
        this.showSuccess(response.message);
        this.invoice = response.invoice || this.invoice;
        this.printing = false;
        
        // Aguarda 2 segundos e volta para lista
        setTimeout(() => {
          this.router.navigate(['/invoices']);
        }, 2000);
      },
      error: (error) => {
        this.showError(error.message);
        this.printing = false;
      }
    });
  }

  /**
   * Retorna o total de unidades da nota
   */
  getTotalQuantity(): number {
    if (!this.invoice) return 0;
    return this.invoice.items.reduce((sum, item) => sum + item.quantity, 0);
  }

  /**
   * Formata data para exibição
   */
  formatDate(dateString: string | undefined): string {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return date.toLocaleString('pt-BR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }

  goBack(): void {
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