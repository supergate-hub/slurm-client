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

## P3: ProxyJump/Bastion SSH 지원
SSH 백엔드에서 ProxyJump/bastion host 경유 연결 지원. HPC 환경에서 bastion을 경유하는 경우가 많음. golang.org/x/crypto/ssh가 ProxyJump를 native로 지원하지 않아 수동 hop chain 구현 필요.
- Effort: M (human) → S (CC)
- Depends on: SSH 백엔드 Phase 1 완료
- Context: Phase 1에서는 직접 SSH 연결만 지원. bastion 필요한 환경은 Phase 2에서 지원.
