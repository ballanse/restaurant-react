﻿using System.Diagnostics.CodeAnalysis;
using Microsoft.Extensions.Logging;
using Restaurant.Server.Api.Models;

namespace Restaurant.Server.Api.Repositories
{
	[ExcludeFromCodeCoverage]
	public class DailyEatingRepository : RepositoryBase<DailyEating>
	{
		public DailyEatingRepository(DatabaseContext context, ILogger<DailyEatingRepository> logger)
			: base(context, logger)
		{
		}
	}
}