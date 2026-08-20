import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { FormGroup, FormControl, Validators, ReactiveFormsModule } from '@angular/forms';
import { InvoiceService, Invoice as InvoiceModel } from '../../services/invoiceService';

@Component({
  selector: 'app-invoice',
  imports: [ReactiveFormsModule],
  templateUrl: './invoice.html',
  styleUrl: './invoice.css',
})
export class Invoice implements OnInit {
  invoices: InvoiceModel[] = [];
  successMessage: string = '';
  errorMessage: string = '';

  invoiceForm = new FormGroup({
    ProductCode: new FormControl('', Validators.required),
    Quantity: new FormControl(1, [Validators.required, Validators.min(1)])
  });

  constructor(
    private invoiceService: InvoiceService, 
    private cdRef: ChangeDetectorRef
  ) {}

  ngOnInit(): void {
    this.loadInvoices();
  }

  loadInvoices(): void {
    this.invoiceService.getInvoices().subscribe({
      next: (data: any) => {
        this.invoices = data.invoices;
        this.cdRef.detectChanges();
      },
      error: (err) => console.error('Erro ao buscar notas:', err)
    });
  }

  onSubmit(): void {
    if (this.invoiceForm.valid) {
      const newInvoice: InvoiceModel = {
        id: 0,
        status: 'PENDING', 
        items: [
          {
            id: 0,
            invoiceId: 0,
            productCode: this.invoiceForm.value.ProductCode ?? '',
            quantity: this.invoiceForm.value.Quantity ?? 1
          }
        ]
      };

      this.invoiceService.createInvoice(newInvoice).subscribe({
        next: () => {
          this.successMessage = 'Nota fiscal gerada!';
          this.invoiceForm.reset({ Quantity: 1 });
          this.loadInvoices();
          setTimeout(() => { this.successMessage = ''; this.cdRef.detectChanges(); }, 3000);
        },
        error: (err) => {
          this.errorMessage = 'Erro ao emitir nota.';
          console.error(err);
          this.cdRef.detectChanges();
        }
      });
    }
  }
}