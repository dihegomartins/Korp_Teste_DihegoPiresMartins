import { Component, OnInit } from '@angular/core'; // Adicionado OnInit
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ProdutoService } from '../../services/produto';

@Component({
  selector: 'app-produto-cadastro',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './produto-cadastro.html',
})
export class ProdutoCadastroComponent implements OnInit {
  // Lista suspensa de produtos
  listaProdutos: any[] = [];

  // Campos para Reposição de Estoque
  idProduto: number | null = null;
  quantidadeParaSomar: number = 0;

  // Campos para Novo Cadastro
  novoProd = {
    codigo: '',
    descricao: '',
    saldo: 0,
  };

  constructor(private produtoService: ProdutoService) {}

  // Carrega os produtos assim que o componente inicia
  ngOnInit(): void {
    this.carregarListaDeProdutos();
  }

  // Busca os produtos no Backend para preencher o Select
  carregarListaDeProdutos() {
    this.produtoService.listarProdutos().subscribe({
      next: (dados) => {
        this.listaProdutos = dados;
      },
      error: (err) => {
        console.error('Erro ao buscar produtos para o select:', err);
      },
    });
  }

  // 1. Função para Somar Saldo (PATCH)
  executarReposicao() {
    if (!this.idProduto || this.quantidadeParaSomar <= 0) {
      alert('Selecione um produto e informe uma quantidade maior que zero.');
      return;
    }

    this.produtoService.adicionarEstoque(this.idProduto, this.quantidadeParaSomar).subscribe({
      next: (res) => {
        alert(`Sucesso: ${res.message}`);
        this.idProduto = null;
        this.quantidadeParaSomar = 0;
        // Atualiza a lista para refletir o novo saldo, se necessário
        this.carregarListaDeProdutos();
      },
      error: (err) => {
        alert('Erro ao repor estoque: ' + (err.error?.error || 'Erro de conexão'));
      },
    });
  }

  // 2. Função para Criar Produto Novo (POST)
  cadastrarNovo() {
    if (!this.novoProd.codigo || !this.novoProd.descricao) {
      alert('Preencha o código e a descrição para o novo cadastro!');
      return;
    }

    this.produtoService.salvarProduto(this.novoProd).subscribe({
      next: (res) => {
        alert('Produto cadastrado com sucesso!');
        this.novoProd = { codigo: '', descricao: '', saldo: 0 };
        // Atualiza a lista suspensa para o novo produto aparecer lá no "Repor Saldo"
        this.carregarListaDeProdutos();
      },
      error: (err) => {
        const msg = err.error?.error || 'Verifique os dados ou a conexão com o servidor';
        alert('Erro ao cadastrar: ' + msg);
      },
    });
  }
}
