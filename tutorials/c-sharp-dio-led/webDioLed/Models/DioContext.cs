namespace webDioLed.Models
{
    using Microsoft.EntityFrameworkCore;
    public class DioContext: DbContext
    {
        public DioContext(DbContextOptions<DioContext> options) : base(options) { }
        public DbSet<Dio> dios { get; set; } = null!;

    }
    
}
