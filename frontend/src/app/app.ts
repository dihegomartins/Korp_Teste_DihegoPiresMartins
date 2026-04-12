import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NotaListaComponent } from './components/nota-lista/nota-lista';
import { ProdutoCadastroComponent } from './components/produto-cadastro/produto-cadastro';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, NotaListaComponent, ProdutoCadastroComponent],
  templateUrl: './app.html',
  styleUrl: './app.scss',
})
export class AppComponent {
  title = 'faturamento-app';
}
