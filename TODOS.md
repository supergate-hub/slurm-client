# TODOS

## ~~P2: MCP 고도화 2단계 — GPU allocation + Queue depth~~ ✅ DONE
Completed in commit c9084b1. `slurm_gpu_allocation` (GRES/TRES per node) + `slurm_queue_depth` (partition/user breakdown).

## ~~P2: RBAC / Audit logging for MCP Server~~ ✅ DONE
Completed in commit c9084b1. 3 access levels (read-only/operator/admin) + JSON audit logging. Config via env vars or clusters.yaml.

## P3: Prometheus/DCGM 연동 데이터 소스
실시간 GPU utilization(nvidia-smi 수준)은 slurmrestd가 제공하지 않음. Prometheus/DCGM exporter에서 데이터 수집하는 별도 어댑터 필요.
- Effort: L (human) → M (CC)
- Depends on: 고도화 2단계 완료 ✅
- Context: MCP 고도화 3단계로 위치. 별도 데이터 소스 연동은 아키텍처 확장 필요.
