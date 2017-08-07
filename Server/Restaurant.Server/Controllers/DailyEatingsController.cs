using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using Restaurant.DataTransferObjects;
using Restaurant.Server.Api.Abstractions.Facades;
using Restaurant.Server.Api.Abstractions.Repositories;
using Restaurant.Server.Api.Constants;
using Restaurant.Server.Api.Models;

namespace Restaurant.Server.Api.Controllers
{
    [Produces("application/json")]
    [Route("/api/[controller]")]
    [AllowAnonymous]
    public class DailyEatingsController : Controller
    {
        private readonly IRepository<DailyEating> _repository;
        private readonly ILogger _logger;
        private readonly IMapperFacade _mapperFacade;

        public DailyEatingsController(
            IRepository<DailyEating> repository,
            ILogger<DailyEatingsController> logger,
            IMapperFacade mapperFacade)
        {
            _repository = repository;
            _logger = logger;
            _mapperFacade = mapperFacade;
        }

        [HttpGet]
        public IEnumerable<DailyEatingDto> Get()
        {
            var dailyEatings = _repository.GetAll();
            return _mapperFacade.Map<IEnumerable<DailyEatingDto>, IQueryable<DailyEating>>(dailyEatings);
        }

        [HttpGet("{id}")]
        public DailyEatingDto Get(Guid id)
        {
            return _mapperFacade.Map<DailyEatingDto, DailyEating>(_repository.Get(id));
        }

        [HttpPost]
        public async Task<IActionResult> Post([FromBody] DailyEatingDto dailyEatingDto, IFormFile receipt = null)
        {
            try
            {
                if (receipt != null)
                {
                    
                    dailyEatingDto.Reciept = receipt.FileName;
                }

                var dailyEating = _mapperFacade.Map<DailyEating, DailyEatingDto>(dailyEatingDto);
                _repository.Create(dailyEating);
                return await _repository.Commit() ? Ok() : (IActionResult)BadRequest();
            }
            catch (Exception)
            {
                return BadRequest();
            }
        }

        [HttpPut("{id}")]
        public async Task<IActionResult> Put(Guid id, [FromBody]DailyEatingDto dto)
        {
            try
            {
                var dailyEating = _mapperFacade.Map<DailyEating, DailyEatingDto>(dto);
                _repository.Update(dailyEating);
                return await _repository.Commit() ? Ok() : (IActionResult)BadRequest();
            }
            catch (Exception)
            {
                return BadRequest();
            }
        }

        [HttpDelete("{id}")]
        public async Task<IActionResult> Delete(Guid id)
        {
            try
            {
                var dailyEating = _repository.Get(id);
                _repository.Delete(dailyEating);
                return await _repository.Commit() ? Ok() : (IActionResult)BadRequest();
            }
            catch (Exception)
            {
                return BadRequest();
            }
        }
    }
}
