﻿using System;
using System.Threading.Tasks;
using ReactiveUI;
using Restaurant.ViewModels;
using Xamarin.Forms.Xaml;

namespace Restaurant.Pages
{
	[XamlCompilation(XamlCompilationOptions.Compile)]
	public partial class FoodDetailPage : FoodDetailPageXaml
	{
		public FoodDetailPage()
		{
			InitializeComponent();
			IsTransparentToolbar = true;
		}

		protected override void OnLoaded()
		{
			BindingContext = ViewModel;

			this.WhenAnyValue(x => x.ViewModel.CurrentOrder.TotalPrice)
				.Subscribe(async totalPrice =>
				{
					var j = totalPrice - 15;
					j = j <= 0 ? 0 : j;
					for (var i = j; i <= totalPrice; i++)
					{
						await Task.Delay(5);
						TotalPrice.Text = $"{i:C}";
					}
				});
		}
	}

	public abstract class FoodDetailPageXaml : BaseContentPage<FoodDetailViewModel>
	{
	}
}