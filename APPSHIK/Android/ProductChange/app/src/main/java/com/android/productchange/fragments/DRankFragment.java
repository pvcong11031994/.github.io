package com.android.productchange.fragments;

import android.app.Dialog;
import android.graphics.Color;
import android.graphics.drawable.ColorDrawable;
import android.os.Bundle;
import android.support.v4.app.DialogFragment;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.view.Window;
import android.widget.Button;
import android.widget.ImageButton;
import android.widget.ImageView;

import com.android.productchange.R;
import com.android.productchange.common.constants.Constants;

/**
 * @author tien-lv
 * @since 2017-12-19
 * Dialog product detail screen
 */

@SuppressWarnings("deprecation")
public class DRankFragment extends DialogFragment implements View.OnClickListener {

    public interface RankDialogListener {
        void onRankSelectedDialog(int itemSelected);
    }

    private int rank;

    ImageView imvArrival, imvPlatform1, imvPlatform2, imvSurface, imvShelter, imvReturn, imvPeroid,
            imvRegular;

    /**
     * Init Dialog Product Detail with custom layout
     */
    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
            Bundle saveInstanceState) {

        // Init custom product detail layout
        View rootView = inflater.inflate(R.layout.layout_rank, container, false);
        if (getDialog().getWindow() != null) {
            getDialog().getWindow().requestFeature(Window.FEATURE_NO_TITLE);
            getDialog().getWindow().setBackgroundDrawable(new ColorDrawable(Color.TRANSPARENT));
        }

        rank = 0;
        if (getArguments() != null) {
            rank = getArguments().getInt(Constants.RANK);
        }

        Button btnArrival = (Button) rootView.findViewById(R.id.btn_arrival);
        Button btnPlatform1 = (Button) rootView.findViewById(R.id.btn_platform1);
        Button btnPlatform2 = (Button) rootView.findViewById(R.id.btn_platform2);
        Button btnSurface = (Button) rootView.findViewById(R.id.btn_surface);
        Button btnShelter = (Button) rootView.findViewById(R.id.btn_shelter);
        Button btnReturn = (Button) rootView.findViewById(R.id.btn_return);
        Button btnPeriod = (Button) rootView.findViewById(R.id.btn_period);
        Button btnRegular = (Button) rootView.findViewById(R.id.btn_regular);
        ImageButton imbBack = (ImageButton) rootView.findViewById(R.id.imb_back);

        imvArrival = (ImageView) rootView.findViewById(R.id.imv_arrival);
        imvPlatform1 = (ImageView) rootView.findViewById(R.id.imv_platform1);
        imvPlatform2 = (ImageView) rootView.findViewById(R.id.imv_platform2);
        imvSurface = (ImageView) rootView.findViewById(R.id.imv_surface);
        imvShelter = (ImageView) rootView.findViewById(R.id.imv_shelter);
        imvReturn = (ImageView) rootView.findViewById(R.id.imv_return);
        imvPeroid = (ImageView) rootView.findViewById(R.id.imv_period);
        imvRegular = (ImageView) rootView.findViewById(R.id.imv_regular);

        loadChecked();

        btnArrival.setOnClickListener(this);
        btnPlatform1.setOnClickListener(this);
        btnPlatform2.setOnClickListener(this);
        btnSurface.setOnClickListener(this);
        btnShelter.setOnClickListener(this);
        btnReturn.setOnClickListener(this);
        btnPeriod.setOnClickListener(this);
        btnRegular.setOnClickListener(this);
        imbBack.setOnClickListener(this);

        return rootView;
    }

    /**
     * Set size layout on start
     */
    @Override
    public void onStart() {
        super.onStart();
        Dialog dialog = getDialog();
        if (dialog != null) {
            if (dialog.getWindow() != null) {
                dialog.getWindow().setLayout(ViewGroup.LayoutParams.FILL_PARENT,
                        ViewGroup.LayoutParams.FILL_PARENT);
            }
        }
    }

    /**
     * load selected
     */
    private void loadChecked() {

        switch (rank) {
            case Constants.RANK_PLATFORM1:
                imvArrival.setVisibility(View.GONE);
                imvPlatform1.setVisibility(View.VISIBLE);
                break;
            case Constants.RANK_PLATFORM2:
                imvArrival.setVisibility(View.GONE);
                imvPlatform2.setVisibility(View.VISIBLE);
                break;
            case Constants.RANK_SURFACE:
                imvArrival.setVisibility(View.GONE);
                imvSurface.setVisibility(View.VISIBLE);
                break;
            case Constants.RANK_SHELDER:
                imvArrival.setVisibility(View.GONE);
                imvShelter.setVisibility(View.VISIBLE);
                break;
            case Constants.RANK_RETURN:
                imvArrival.setVisibility(View.GONE);
                imvReturn.setVisibility(View.VISIBLE);
                break;
            case Constants.RANK_PERIOD:
                imvArrival.setVisibility(View.GONE);
                imvPeroid.setVisibility(View.VISIBLE);
                break;
            case Constants.RANK_REGULAR:
                imvArrival.setVisibility(View.GONE);
                imvRegular.setVisibility(View.VISIBLE);
                break;
            default:
                break;
        }
    }

    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.btn_platform1:
                rank = Constants.RANK_PLATFORM1;
                break;
            case R.id.btn_platform2:
                rank = Constants.RANK_PLATFORM2;
                break;
            case R.id.btn_surface:
                rank = Constants.RANK_SURFACE;
                break;
            case R.id.btn_shelter:
                rank = Constants.RANK_SHELDER;
                break;
            case R.id.btn_return:
                rank = Constants.RANK_RETURN;
                break;
            case R.id.btn_arrival:
                rank = Constants.RANK_ARRIVAL;
                break;
            case R.id.btn_period:
                rank = Constants.RANK_PERIOD;
                break;
            case R.id.btn_regular:
                rank = Constants.RANK_REGULAR;
                break;
            default:
                break;
        }
        RankDialogListener activity = (RankDialogListener) getActivity();
        activity.onRankSelectedDialog(rank);
        dismiss();
    }
}
