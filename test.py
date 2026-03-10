import time
from concurrent.futures import ProcessPoolExecutor
import multiprocessing
import os

def cpu_bound(number):
    """CPU密集型任务：计算平方和"""
    result = sum(i * i for i in range(number))
    # 使用 os.getpid() 显示是哪个进程在运行
    print(f"[PID {os.getpid()}] 计算 {number:,} 的结果: {result}")
    return result

def main():
    numbers = [10000000 + x for x in range(20)]
    
    print(f"CPU 核心数: {multiprocessing.cpu_count()}")
    print(f"任务数: {len(numbers)} 个")
    print("=" * 50)
    
    # 1. 单进程版本（顺序执行）
    print("【单进程版本】")
    start = time.perf_counter()
    results_single = [cpu_bound(n) for n in numbers]
    single_duration = time.perf_counter() - start
    print(f"⏱️  单进程耗时: {single_duration:.2f} 秒\n")
    
    # 2. 多进程版本（并行执行）
    print("【多进程版本】")
    start = time.perf_counter()
    
    # ProcessPoolExecutor 会自动使用所有 CPU 核心
    with ProcessPoolExecutor() as executor:
        # map 保持输入输出顺序一致
        results_multi = list(executor.map(cpu_bound, numbers))
    
    multi_duration = time.perf_counter() - start
    print(f"⏱️  多进程耗时: {multi_duration:.2f} 秒\n")
    
    # 3. 性能对比
    print("=" * 50)
    print("【性能分析】")
    print(f"单进程耗时: {single_duration:.2f} 秒")
    print(f"多进程耗时: {multi_duration:.2f} 秒")
    print(f"加速比: {single_duration/multi_duration:.2f}x")
    print(f"效率提升: {(1 - multi_duration/single_duration)*100:.1f}%")
    
    # 验证结果一致性
    assert results_single == results_multi, "计算结果不一致！"
    print("\n✅ 结果验证通过")

if __name__ == '__main__':
    main()