import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError, BehaviorSubject } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';
import { Invoice, CreateInvoiceDTO, PrintResponse } from '../models/invoice.model';

@Injectable({
  providedIn: 'root'
})
export class InvoiceService {
  private apiUrl = 'http://localhost:8082/api/invoices';
  
  // BehaviorSubject para manter lista de notas em memória
  private invoicesSubject = new BehaviorSubject<Invoice[]>([]);
  public invoices$ = this.invoicesSubject.asObservable();

  // Subject para indicar quando está imprimindo (loading)
  private printingSubject = new BehaviorSubject<boolean>(false);
  public printing$ = this.printingSubject.asObservable();

  constructor(private http: HttpClient) {}

  /**
   * Lista todas as notas fiscais
   */
  getInvoices(): Observable<Invoice[]> {
    return this.http.get<Invoice[]>(this.apiUrl).pipe(
      tap(invoices => this.invoicesSubject.next(invoices)),
      catchError(this.handleError)
    );
  }

  /**
   * Busca uma nota fiscal por ID
   */
  getInvoice(id: string): Observable<Invoice> {
    return this.http.get<Invoice>(`${this.apiUrl}/${id}`).pipe(
      catchError(this.handleError)
    );
  }

  /**
   * Cria uma nova nota fiscal
   */
  createInvoice(invoice: CreateInvoiceDTO): Observable<Invoice> {
    return this.http.post<Invoice>(this.apiUrl, invoice).pipe(
      tap(() => this.getInvoices().subscribe()), // Recarrega lista
      catchError(this.handleError)
    );
  }

  /**
   * Imprime (fecha) uma nota fiscal
   * Esta é a operação mais crítica do sistema
   */
  printInvoice(id: string): Observable<PrintResponse> {
    this.printingSubject.next(true); // Ativa indicador de loading
    
    return this.http.post<PrintResponse>(`${this.apiUrl}/${id}/print`, {}).pipe(
      tap(() => {
        this.printingSubject.next(false); // Desativa loading
        this.getInvoices().subscribe(); // Recarrega lista
      }),
      catchError(error => {
        this.printingSubject.next(false); // Desativa loading mesmo em erro
        return this.handleError(error);
      })
    );
  }

  /**
   * Tratamento centralizado de erros
   */
  private handleError(error: HttpErrorResponse) {
    let errorMessage = 'Ocorreu um erro desconhecido';
    
    if (error.error instanceof ErrorEvent) {
      // Erro do lado do cliente
      errorMessage = `Erro: ${error.error.message}`;
    } else {
      // Erro do lado do servidor
      if (error.status === 0) {
        errorMessage = 'Não foi possível conectar ao servidor. Verifique se o Billing Service está rodando.';
      } else if (error.status === 404) {
        errorMessage = 'Nota fiscal não encontrada';
      } else if (error.status === 400) {
        errorMessage = error.error?.error || 'Dados inválidos';
      } else if (error.status === 503) {
        // Falha na comunicação com Stock Service
        errorMessage = 'Falha ao comunicar com o serviço de estoque. Tente novamente.';
      } else if (error.error?.error) {
        errorMessage = error.error.error;
      } else if (error.error?.message) {
        errorMessage = error.error.message;
      } else {
        errorMessage = `Erro do servidor: ${error.status}`;
      }
    }
    
    console.error('Erro no InvoiceService:', error);
    return throwError(() => new Error(errorMessage));
  }
}