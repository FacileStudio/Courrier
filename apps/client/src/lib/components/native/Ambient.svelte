<div class="ambient" aria-hidden="true">
	<span class="orb orb-a"></span>
	<span class="orb orb-b"></span>
</div>

<style>
	.ambient {
		position: absolute;
		inset: 0;
		overflow: hidden;
		pointer-events: none;
	}

	.orb {
		position: absolute;
		aspect-ratio: 1;
		border-radius: 9999px;
		filter: blur(90px);
		opacity: 0.12;
		background-color: var(--color-fc-accent);
		animation: drift 26s var(--ease-fc) infinite alternate;
	}

	.orb-a {
		width: min(70vw, 520px);
		top: -14%;
		left: -16%;
	}

	.orb-b {
		width: min(60vw, 440px);
		bottom: -16%;
		right: -14%;
		opacity: 0.08;
		animation-duration: 32s;
		animation-delay: -8s;
	}

	@keyframes drift {
		from {
			transform: translate3d(0, 0, 0) scale(1);
		}
		to {
			transform: translate3d(6vw, 5vh, 0) scale(1.12);
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.orb {
			animation: none;
		}
	}
</style>
