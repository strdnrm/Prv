namespace task5
{
    public partial class Form1 : Form
    {
        int maxIterations = 100;
        public Form1()
        {
            InitializeComponent();
        }
        private async void Form_Load(object sender, EventArgs e)
        {
            await GenerateMandelbrotSetAsync();
        }

        private async Task GenerateMandelbrotSetAsync()
        {
            int width = pictureBox1.Width;
            int height = pictureBox1.Height;
        
            Bitmap bitmap = new Bitmap(width, height);
            Mutex mutexObj = new();
            await Task.Run(() =>
            {
                Parallel.For(0, height, y =>
                {
                    for (int x = 0; x < width; x++)
                    {
                        double a = Map(x, 0, width, -2.5, 1); 
                        double b = Map(y, 0, height, -1, 1); 

                        int iteration = CalculateMandelbrot(a, b);

                        int rCol = (iteration * 10) % 256;
                        int gCol = (iteration * 5) % 256; 
                        int bCol = (iteration * 2) % 256; 
                        Color color = iteration == maxIterations ? Color.Black : Color.FromArgb(rCol, gCol, bCol);
                        mutexObj.WaitOne();
                        bitmap.SetPixel(x, y, color);
                        mutexObj.ReleaseMutex();
                    }
                });
            });

            pictureBox1.Image = bitmap;
        }

        private int CalculateMandelbrot(double a, double b)
        {
            double x = 0;
            double y = 0;
            int iteration = 0;

            while (x * x + y * y <= 4 && iteration < maxIterations)
            {
                double xtemp = x * x - y * y + a;
                y = 2 * x * y + b;
                x = xtemp;
                iteration++;
            }

            return iteration;
        }

        private double Map(double value, double start1, double stop1, double start2, double stop2)
        {
            return start2 + (stop2 - start2) * ((value - start1) / (stop1 - start1));
        }
    }
}