using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using webDioLed.Models;

namespace webDioLed.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class DiosController : ControllerBase
    {
        private readonly DioContext _context;

        public DiosController(DioContext context)
        {
            _context = context;
        }

        // GET: api/Dios
        [HttpGet]
        public async Task<ActionResult<IEnumerable<Dio>>> Getdios()
        {
          if (_context.dios == null)
          {
              return NotFound();
          }
            return await _context.dios.ToListAsync();
        }

        // GET: api/Dios/5
        [HttpGet("{id}")]
        public async Task<ActionResult<Dio>> GetDio(int id)
        {
          if (_context.dios == null)
          {
              return NotFound();
          }
            var dio = await _context.dios.FindAsync(id);

            if (dio == null)
            {
                return NotFound();
            }

            return dio;
        }

        // PUT: api/Dios/5
        // To protect from overposting attacks, see https://go.microsoft.com/fwlink/?linkid=2123754
        [HttpPut("{id}")]
        public async Task<IActionResult> PutDio(int id, Dio dio)
        {
            if (id != dio.Id)
            {
                return BadRequest();
            }

            _context.Entry(dio).State = EntityState.Modified;

            try
            {
                await _context.SaveChangesAsync();
            }
            catch (DbUpdateConcurrencyException)
            {
                if (!DioExists(id))
                {
                    return NotFound();
                }
                else
                {
                    throw;
                }
            }

            return NoContent();
        }

        // POST: api/Dios
        // To protect from overposting attacks, see https://go.microsoft.com/fwlink/?linkid=2123754
        [HttpPost]
        public async Task<ActionResult<Dio>> PostDio(Dio dio)
        {
          if (_context.dios == null)
          {
              return Problem("Entity set 'DioContext.dios'  is null.");
          }
            _context.dios.Add(dio);
            await _context.SaveChangesAsync();

            return CreatedAtAction("GetDio", new { id = dio.Id }, dio);
        }

        // DELETE: api/Dios/5
        [HttpDelete("{id}")]
        public async Task<IActionResult> DeleteDio(int id)
        {
            if (_context.dios == null)
            {
                return NotFound();
            }
            var dio = await _context.dios.FindAsync(id);
            if (dio == null)
            {
                return NotFound();
            }

            _context.dios.Remove(dio);
            await _context.SaveChangesAsync();

            return NoContent();
        }

        private bool DioExists(int id)
        {
            return (_context.dios?.Any(e => e.Id == id)).GetValueOrDefault();
        }
    }
}
